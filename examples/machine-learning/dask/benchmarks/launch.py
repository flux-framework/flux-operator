import argparse
import operator
import time
import math
from dask_jobqueue import FluxCluster
from distributed import Client

import pandas as pd
from dask.distributed import Client, wait
import dask.config

# Being lazy, really. Don't do this ;)
client = None

# These are tasks we want to test! This code is from:
# https://gist.github.com/mrocklin/4c198b13e92f881161ef175810c7f6bc#file-scaling-gcs-ipynb
# we are trying to reproduce it, but running in the flux operator

def slowinc(x, delay=0.1):
    time.sleep(delay)
    return x + 1

def slowadd(x, y, delay=0.1):
    time.sleep(delay)
    return x + y

def slowsum(L, delay=0.1):
    time.sleep(delay)
    return sum(L)

def inc(x):
    return x + 1

def tasks(n):
    yield 'task map fast tasks', 'tasks', n * 200
    
    futures = client.map(inc, range(n * 200))
    wait(futures)
    
    yield 'task map 100ms tasks', 'tasks', n * 100

    futures = client.map(slowinc, range(100 * n))
    wait(futures)
        
    yield 'task map 1s tasks', 'tasks', n * 4

    futures = client.map(slowinc, range(4 * n), delay=1)
    wait(futures)

    yield 'tree reduction fast tasks', 'tasks', 2**7 * n
    
    from dask import delayed

    L = range(2**7 * n)
    while len(L) > 1:
        L = list(map(delayed(operator.add), L[0::2], L[1::2]))

    L[0].compute()
    
    yield 'tree reduction 100ms tasks', 'tasks', 2**6 * n * 2
    
    from dask import delayed

    L = range(2**6 * n)
    while len(L) > 1:
        L = list(map(delayed(slowadd), L[0::2], L[1::2]))

    L[0].compute()
    
    yield 'sequential', 'tasks', 100

    x = 1

    for i in range(100):
        x = delayed(inc)(x)
        
    x.compute()
    
    yield 'dynamic tree reduction fast tasks', 'tasks', 100 * n
    
    from dask.distributed import as_completed
    futures = client.map(inc, range(n * 100))
    
    pool = as_completed(futures)
    batches = pool.batches()
    
    while True:
        try:
            batch = next(batches)
            if len(batch) == 1:
                batch += next(batches)
        except StopIteration:
            break
        future = client.submit(sum, batch)
        pool.add(future)
        
    yield 'dynamic tree reduction 100ms tasks', 'tasks', 100 * n
    
    from dask.distributed import as_completed
    futures = client.map(slowinc, range(n * 20))
    
    pool = as_completed(futures)
    batches = pool.batches()
    
    while True:
        try:
            batch = next(batches)
            if len(batch) == 1:
                batch += next(batches)
        except StopIteration:
            break
        future = client.submit(slowsum, batch)
        pool.add(future)

        
    yield 'nearest neighbor fast tasks', 'tasks', 100 * n * 2
    
    L = range(100 * n)
    L = client.map(operator.add, L[:-1], L[1:])
    L = client.map(operator.add, L[:-1], L[1:])
    wait(L)
    
    yield 'nearest neighbor 100ms tasks', 'tasks', 20 * n * 2
    
    L = range(20 * n)
    L = client.map(slowadd, L[:-1], L[1:])
    L = client.map(slowadd, L[:-1], L[1:])
    wait(L)

def arrays(n):
    import dask.array as da
    N = int(5000 * math.sqrt(n))
    x = da.random.randint(0, 10000, size=(N, N), chunks=(2000, 2000))
    
    yield 'create random', 'MB', x.nbytes / 1e6
    
    x = x.persist()
    wait(x)
    
    yield 'blockwise 100ms tasks', 'MB', x.nbytes / 1e6
    
    y = x.map_blocks(slowinc, dtype=x.dtype).persist()
    wait(y)
    
    yield 'random access', 'bytes', 8
    
    x[1234, 4567].compute()
   
    yield 'reduction', 'MB', x.nbytes / 1e6
    
    x.std().compute()
    
    yield 'reduction along axis', 'MB', x.nbytes / 1e6
    
    x.std(axis=0).compute()
    
    yield 'elementwise computation', 'MB', x.nbytes / 1e6
    
    y = da.sin(x) ** 2 + da.cos(x) ** 2
    y = y.persist()
    wait(y)    
    
    yield 'rechunk small', 'MB', x.nbytes / 1e6
    
    y = x.rechunk((20000, 200)).persist()
    wait(y)
    
    yield 'rechunk large', 'MB', x.nbytes / 1e6
    
    y = y.rechunk((200, 20000)).persist()
    wait(y)
    
    yield 'transpose addition', 'MB', x.nbytes / 1e6
    y = x + x.T
    y = y.persist()
    wait(y)
    
    yield 'nearest neighbor fast tasks', 'MB', x.nbytes / 1e6
    
    y = x.map_overlap(inc, depth=1).persist()
    wait(y)   
        
    yield 'nearest neighbor 100ms tasks', 'MB', x.nbytes / 1e6
    
    y = x.map_overlap(slowinc, depth=1, delay=0.1).persist()
    wait(y)    

def dataframes(n):
    import dask.array as da
    import dask.dataframe as dd
    N = 2000000 * n
    
    x = da.random.randint(0, 10000, size=(N, 10), chunks=(1000000, 10))

    
    yield 'create random', 'MB', x.nbytes / 1e6
    
    df = dd.from_dask_array(x).persist()
    wait(df)
    
    yield 'blockwise 100ms tasks', 'MB', x.nbytes / 1e6
    
    wait(df.map_partitions(slowinc, meta=df).persist())
    
    yield 'arithmetic', 'MB', x.nbytes / 1e6
    
    y = (df[0] + 1 + 2 + 3 + 4 + 5 + 6 + 7 + 8 + 9 + 10).persist()
    wait(y)
    
    yield 'random access', 'bytes', 8
    
    df.loc[123456].compute()
    
    yield 'dataframe reduction', 'MB', x.nbytes / 1e6
    
    df.std().compute()
    
    yield 'series reduction', 'MB', x.nbytes / 1e6 / 10
    
    df[3].std().compute()
    
    yield 'groupby reduction', 'MB', x.nbytes / 1e6
    
    df.groupby(0)[1].mean().compute()
    
    yield 'groupby apply (full shuffle)', 'MB', x.nbytes / 1e6
    
    df.groupby(0).apply(len).compute()
    
    yield 'set index (full shuffle)', 'MB', x.nbytes / 1e6
    
    wait(df.set_index(1).persist())
    
    yield 'rolling aggregations', 'MB', x.nbytes / 1e6
    
    wait(df.rolling(5).mean().persist())



def run(func, client=None):
    client.restart()
    n = sum(client.ncores().values())
    coroutine = func(n)

    name, unit, numerator = next(coroutine)
    out = []
    while True:
        time.sleep(1)
        start = time.time()
        try:
            next_name, next_unit, next_numerator = next(coroutine)
        except StopIteration:
            break
        finally:
            end = time.time()
            record = {'name': name, 
                      'duration': end - start, 
                      'unit': unit + '/s', 
                      'rate': numerator / (end - start), 
                      'n': n,
                      'collection': func.__name__}
            out.append(record)
        name = next_name
        unit = next_unit
        numerator = next_numerator
    return pd.DataFrame(out)


# Ensure the place dask writes files is shared by all nodes!
dask.config.set(temporary_directory='/tmp/workflow/tmp')

def get_parser():
    parser = argparse.ArgumentParser(
        description="Flux Basic Experiment Runner",
        formatter_class=argparse.RawTextHelpFormatter,
    )
    parser.add_argument("--iters", help="iterations to run", type=int, default=3)
    parser.add_argument("--workers", help="max number of worker nodes in the cluster", type=int, default=4)
    parser.add_argument("--cores", help="cores to give to flux cluster spec", type=int, default=1)
    parser.add_argument('--memory', help='dummy memory variable (is not used)', default="2GB")
    parser.add_argument('--timeout', help='default timeout to wait for workers (60 seconds)', type=int, default=60)
    return parser


def main():
    parser = get_parser()

    # If an error occurs while parsing the arguments, the interpreter will exit with value 2
    args, _ = parser.parse_known_args()

    # Show args to the user
    print("       iters: %s" % args.iters)
    print("       nodes: %s" % args.workers)
    print("       cores: %s" % args.cores)
    print("     timeout: %s" % args.timeout)

    with FluxCluster(cores=args.cores, processes=1, memory=args.memory) as cluster:
        cluster.adapt()

        # Don't use global variables in actual software
        global client 

        with Client(cluster) as client:
            print(f'Preparing cluster with {args.workers} workers')
            cluster.scale(args.workers)
            client.wait_for_workers(args.workers, timeout=args.timeout)            
            print(f'Client: {client}')

            # https://gist.github.com/mrocklin/4c198b13e92f881161ef175810c7f6bc#file-scaling-gcs-ipynb
            # This value of n gets the cores based on the size we've selected!
            n = sum(client.ncores().values())
            print(f'Value of n is {n}')

            L = []
            for i in range(args.iters):
                for func in [tasks, arrays, dataframes]:
                    print(f"Iteration {i}", func.__name__)
                    df = run(func, client=client)
                    L.append(df)

            # Assemble into one data frame to save for further analysis
            print(f'Experiments are done! Saving data to dask-experiments-{args.workers}-raw.csv...')
            ddf = pd.concat(L)
            ddf.to_csv(f"dask-experiments-{args.workers}-raw.csv", index=False)

if __name__ == "__main__":
    main()

