import argparse
from time import sleep
from dask_jobqueue import FluxCluster
from distributed import Client

import numpy as np
from dask.distributed import Client
import dask.config

import joblib
from sklearn.datasets import load_digits
from sklearn.model_selection import RandomizedSearchCV
from sklearn.svm import SVC


# Ensure the place dask writes files is shared by all nodes!
dask.config.set(temporary_directory='/tmp/workflow/tmp')

def get_parser():
    parser = argparse.ArgumentParser(
        description="Flux Basic Experiment Runner",
        formatter_class=argparse.RawTextHelpFormatter,
    )
    parser.add_argument("--workers", help="number of worker nodes in the cluster", type=int, default=4)
    parser.add_argument("--cores", help="cores to give to flux cluster spec", type=int, default=1)
    parser.add_argument('--memory', help='dummy memory variable (is not used)', default="2GB")
    parser.add_argument('--timeout', help='default timeout to wait for workers (60 seconds)', type=int, default=60)
    return parser


def main():
    parser = get_parser()

    # If an error occurs while parsing the arguments, the interpreter will exit with value 2
    args, _ = parser.parse_known_args()

    # Show args to the user
    print("     nodes: %s" % args.workers)
    print("     cores: %s" % args.cores)
    print("   timeout: %s" % args.timeout)

    # For flux, memory doesn't matter (it's ignored)
    with FluxCluster(cores=args.cores, processes=1, memory=args.memory) as cluster:
        cluster.adapt()
        with Client(cluster) as client:
            cluster.scale(args.workers)
            client.wait_for_workers(args.workers, timeout=args.timeout)
            digits = load_digits()

            param_space = {
                'C': np.logspace(-6, 6, 13),
                'gamma': np.logspace(-8, 8, 17),
                'tol': np.logspace(-4, -1, 4),
                'class_weight': [None, 'balanced'],
            }
            model = SVC(kernel='rbf')
            search = RandomizedSearchCV(model, param_space, cv=3, n_iter=50, verbose=10)

            print('Fitting model...')
            with joblib.parallel_backend('dask'):
                search.fit(digits.data, digits.target)

            print('CV Results:')
            print(search.cv_results_)
            print()
            print('Best Params')
            print(search.best_params_)
            print('Sleeping for two minutes to keep job alive if you want to interact!')
            sleep(120)

if __name__ == "__main__":
    main()

