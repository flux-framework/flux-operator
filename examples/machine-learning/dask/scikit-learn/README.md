# Dask with scikit-learn and Flux!

This will test running [Dask and scikit-learn](https://ml.dask.org/joblib.html) with Flux.
Our goal is to do this simply, and eventually extend this to a more complex hierarchy of jobs.

## Usage

First, let's create a kind cluster.

```bash
$ kind create cluster --config ../../../kind-config.yaml
```

And then install the operator, create the namespace, and apply the MiniCluster YAML here.

```bash
$ kubectl apply -f ../../../dist/flux-operator.yaml
$ kubectl create namespace flux-operator
$ kubectl apply -f ./minicluster.yaml
```

This will install dependencies (dask and scikit-learn) directly into the base image, and then mount
the current directory (in the pods as `/tmp/workflow`). We run the [launch.py](launch.py)
script from the broker, and you can inspect this script to see how we create and connect
to Flux. There is currently a generous timeout (60 seconds) to ensure
the cluster is ready. You can watch logs doing the following: 

```bash
$ kubectl logs -n flux-operator flux-sample-0-7tx7s -f
```

Since I wrote this as a demo, I basically did model training, and then printed the final CV results
and parameters. You'll see this in the log (expand the section below):

<details>

<summary>Output of Scikit-Learn Example</summary>

```console
...
broker.info[0]: rc1.0: running /etc/flux/rc1.d/01-sched-fluxion
sched-fluxion-resource.info[0]: version 214aa27
sched-fluxion-resource.warning[0]: create_reader: allowlist unsupported
sched-fluxion-resource.info[0]: populate_resource_db: loaded resources from core's resource.acquire
sched-fluxion-qmanager.info[0]: version 214aa27
broker.info[0]: rc1.0: running /etc/flux/rc1.d/02-cron
broker.info[0]: rc1.0: /etc/flux/rc1 Exited (rc=0) 3.8s
broker.info[0]: rc1-success: init->quorum 3.8019s
broker.info[0]: online: flux-sample-0 (ranks 0)
broker.info[0]: online: flux-sample-[0-3] (ranks 0-3)
broker.info[0]: quorum-full: quorum->run 3.99891s
     nodes: 4
     cores: 2
   timeout: 60
Fitting model...
Fitting 3 folds for each of 50 candidates, totalling 150 fits
CV Results:
{'mean_fit_time': array([4.08554602, 5.423575  , 3.39528402, 2.85911171, 4.831592  ,
       3.88271205, 3.02563477, 2.47843599, 0.68565289, 2.93963766,
       3.39221478, 3.79691935, 3.96530684, 3.82564553, 4.97321709,
       3.84288883, 1.15530101, 1.35951893, 3.7286648 , 4.03004575,
       3.66520365, 3.37004209, 3.13437613, 2.99476473, 2.26174323,
       3.39319603, 2.89181201, 2.62760369, 0.65904252, 0.39908274,
       3.00021768, 0.34174371, 3.26659242, 3.26887321, 3.08342997,
       2.61445681, 3.21243397, 3.22617459, 3.29642828, 2.41557391,
       3.36702911, 3.24384109, 2.87103478, 2.92750065, 0.34344927,
       0.5176696 , 2.61106515, 0.41769854, 0.67651065, 0.86017307]), 'std_fit_time': array([1.29295808, 0.44993044, 1.13105463, 0.24890695, 1.53956746,
       0.93328344, 0.66994432, 0.16029377, 0.17592682, 0.12820939,
       0.37487107, 0.68207106, 0.36123274, 0.27662093, 0.44998325,
       0.11232528, 0.15331819, 0.29612327, 0.23260674, 0.20277183,
       0.10520431, 0.34012566, 0.15273326, 0.25142427, 0.26527932,
       0.10249256, 0.11939426, 0.40004664, 0.02033139, 0.06673623,
       0.05008267, 0.05227693, 0.29535088, 0.14343989, 0.26508654,
       0.20566045, 0.40562526, 0.30273003, 0.62674814, 0.09708462,
       0.41936624, 0.34741114, 0.15964411, 0.2086272 , 0.04933951,
       0.01189969, 0.22959803, 0.08866899, 0.05076011, 0.12571382]), 'mean_score_time': array([2.26498397, 1.77561784, 1.6284159 , 2.33371989, 1.4934899 ,
       1.65991545, 1.01165613, 1.26928711, 0.80194203, 1.01934727,
       1.14952954, 1.31477968, 2.05600015, 1.86737108, 1.60693971,
       1.95683328, 0.96571255, 1.25838002, 0.99948716, 1.46727578,
       0.9671642 , 1.23076653, 1.0414288 , 1.23779511, 1.16715463,
       1.123487  , 0.96316584, 1.07391651, 0.62996387, 0.42405613,
       1.30649598, 0.31687458, 1.05327328, 1.01003925, 1.0420777 ,
       0.98385247, 1.26904154, 1.24301561, 1.3658274 , 1.24738073,
       1.36005807, 1.32252598, 1.04261764, 1.03164657, 0.29092725,
       0.48727504, 0.32448387, 0.60944939, 0.42013979, 0.33401561]), 'std_score_time': array([0.20775316, 0.24286407, 0.11926821, 0.40169773, 0.55680183,
       0.78485425, 0.17111927, 0.06430806, 0.1293463 , 0.23254882,
       0.26622301, 0.2448651 , 0.21866524, 0.3043432 , 0.11532111,
       0.35099115, 0.02410389, 0.02678309, 0.2368529 , 0.13744223,
       0.1707387 , 0.17363425, 0.07558077, 0.02828425, 0.04064274,
       0.07727504, 0.03594565, 0.03109774, 0.16673041, 0.12429896,
       0.26720998, 0.06015071, 0.24929352, 0.15515904, 0.10645853,
       0.12264939, 0.19762409, 0.29526357, 0.08137586, 0.23999344,
       0.21649012, 0.16955293, 0.10475885, 0.22040064, 0.07748718,
       0.02336777, 0.0849569 , 0.20914747, 0.09983006, 0.12743123]), 'param_tol': masked_array(data=[0.0001, 0.0001, 0.01, 0.1, 0.0001, 0.1, 0.001, 0.1,
                   0.1, 0.1, 0.01, 0.0001, 0.01, 0.1, 0.0001, 0.1, 0.0001,
                   0.01, 0.0001, 0.1, 0.01, 0.1, 0.001, 0.1, 0.01, 0.0001,
                   0.01, 0.001, 0.01, 0.01, 0.0001, 0.1, 0.01, 0.1, 0.001,
                   0.001, 0.1, 0.01, 0.0001, 0.01, 0.001, 0.01, 0.1,
                   0.0001, 0.01, 0.0001, 0.1, 0.01, 0.0001, 0.001],
             mask=[False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False],
       fill_value='?',
            dtype=object), 'param_gamma': masked_array(data=[10000.0, 1000.0, 100000.0, 0.01, 100000.0, 0.1, 1e-08,
                   1000.0, 1e-05, 10.0, 0.0001, 0.1, 0.01, 1e-08,
                   100000000.0, 10.0, 0.0001, 0.001, 1e-08, 0.001, 1e-07,
                   100000000.0, 10000000.0, 100000.0, 100000.0, 1e-05,
                   100.0, 100000.0, 1e-07, 0.0001, 100000.0, 1e-05,
                   100000000.0, 1.0, 1.0, 10000.0, 0.1, 1.0, 10.0,
                   100000.0, 0.01, 100.0, 100000.0, 100000.0, 1e-05,
                   0.0001, 100.0, 0.0001, 1e-05, 1e-07],
             mask=[False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False],
       fill_value='?',
            dtype=object), 'param_class_weight': masked_array(data=['balanced', None, 'balanced', 'balanced', 'balanced',
                   'balanced', None, None, 'balanced', 'balanced',
                   'balanced', None, 'balanced', None, 'balanced', None,
                   None, None, None, 'balanced', 'balanced', 'balanced',
                   None, 'balanced', 'balanced', 'balanced', None, None,
                   None, None, None, None, None, 'balanced', 'balanced',
                   'balanced', None, 'balanced', None, None, 'balanced',
                   None, None, 'balanced', None, 'balanced', None, None,
                   'balanced', 'balanced'],
             mask=[False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False],
       fill_value='?',
            dtype=object), 'param_C': masked_array(data=[100.0, 1000.0, 1e-05, 0.001, 1000.0, 0.01, 1e-05,
                   10000.0, 10.0, 1000.0, 0.0001, 0.1, 100000.0, 1.0,
                   1000000.0, 1000.0, 1.0, 1000.0, 0.1, 0.001, 0.1, 0.01,
                   1.0, 0.1, 1.0, 0.0001, 1000.0, 0.01, 1000.0, 1000000.0,
                   1.0, 100.0, 100.0, 1e-06, 1000.0, 0.01, 0.0001, 0.1,
                   10000.0, 0.01, 0.0001, 100000.0, 0.01, 1e-06, 10000.0,
                   1000000.0, 1000000.0, 10.0, 1000.0, 1000.0],
             mask=[False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False],
       fill_value='?',
            dtype=object), 'params': [{'tol': 0.0001, 'gamma': 10000.0, 'class_weight': 'balanced', 'C': 100.0}, {'tol': 0.0001, 'gamma': 1000.0, 'class_weight': None, 'C': 1000.0}, {'tol': 0.01, 'gamma': 100000.0, 'class_weight': 'balanced', 'C': 1e-05}, {'tol': 0.1, 'gamma': 0.01, 'class_weight': 'balanced', 'C': 0.001}, {'tol': 0.0001, 'gamma': 100000.0, 'class_weight': 'balanced', 'C': 1000.0}, {'tol': 0.1, 'gamma': 0.1, 'class_weight': 'balanced', 'C': 0.01}, {'tol': 0.001, 'gamma': 1e-08, 'class_weight': None, 'C': 1e-05}, {'tol': 0.1, 'gamma': 1000.0, 'class_weight': None, 'C': 10000.0}, {'tol': 0.1, 'gamma': 1e-05, 'class_weight': 'balanced', 'C': 10.0}, {'tol': 0.1, 'gamma': 10.0, 'class_weight': 'balanced', 'C': 1000.0}, {'tol': 0.01, 'gamma': 0.0001, 'class_weight': 'balanced', 'C': 0.0001}, {'tol': 0.0001, 'gamma': 0.1, 'class_weight': None, 'C': 0.1}, {'tol': 0.01, 'gamma': 0.01, 'class_weight': 'balanced', 'C': 100000.0}, {'tol': 0.1, 'gamma': 1e-08, 'class_weight': None, 'C': 1.0}, {'tol': 0.0001, 'gamma': 100000000.0, 'class_weight': 'balanced', 'C': 1000000.0}, {'tol': 0.1, 'gamma': 10.0, 'class_weight': None, 'C': 1000.0}, {'tol': 0.0001, 'gamma': 0.0001, 'class_weight': None, 'C': 1.0}, {'tol': 0.01, 'gamma': 0.001, 'class_weight': None, 'C': 1000.0}, {'tol': 0.0001, 'gamma': 1e-08, 'class_weight': None, 'C': 0.1}, {'tol': 0.1, 'gamma': 0.001, 'class_weight': 'balanced', 'C': 0.001}, {'tol': 0.01, 'gamma': 1e-07, 'class_weight': 'balanced', 'C': 0.1}, {'tol': 0.1, 'gamma': 100000000.0, 'class_weight': 'balanced', 'C': 0.01}, {'tol': 0.001, 'gamma': 10000000.0, 'class_weight': None, 'C': 1.0}, {'tol': 0.1, 'gamma': 100000.0, 'class_weight': 'balanced', 'C': 0.1}, {'tol': 0.01, 'gamma': 100000.0, 'class_weight': 'balanced', 'C': 1.0}, {'tol': 0.0001, 'gamma': 1e-05, 'class_weight': 'balanced', 'C': 0.0001}, {'tol': 0.01, 'gamma': 100.0, 'class_weight': None, 'C': 1000.0}, {'tol': 0.001, 'gamma': 100000.0, 'class_weight': None, 'C': 0.01}, {'tol': 0.01, 'gamma': 1e-07, 'class_weight': None, 'C': 1000.0}, {'tol': 0.01, 'gamma': 0.0001, 'class_weight': None, 'C': 1000000.0}, {'tol': 0.0001, 'gamma': 100000.0, 'class_weight': None, 'C': 1.0}, {'tol': 0.1, 'gamma': 1e-05, 'class_weight': None, 'C': 100.0}, {'tol': 0.01, 'gamma': 100000000.0, 'class_weight': None, 'C': 100.0}, {'tol': 0.1, 'gamma': 1.0, 'class_weight': 'balanced', 'C': 1e-06}, {'tol': 0.001, 'gamma': 1.0, 'class_weight': 'balanced', 'C': 1000.0}, {'tol': 0.001, 'gamma': 10000.0, 'class_weight': 'balanced', 'C': 0.01}, {'tol': 0.1, 'gamma': 0.1, 'class_weight': None, 'C': 0.0001}, {'tol': 0.01, 'gamma': 1.0, 'class_weight': 'balanced', 'C': 0.1}, {'tol': 0.0001, 'gamma': 10.0, 'class_weight': None, 'C': 10000.0}, {'tol': 0.01, 'gamma': 100000.0, 'class_weight': None, 'C': 0.01}, {'tol': 0.001, 'gamma': 0.01, 'class_weight': 'balanced', 'C': 0.0001}, {'tol': 0.01, 'gamma': 100.0, 'class_weight': None, 'C': 100000.0}, {'tol': 0.1, 'gamma': 100000.0, 'class_weight': None, 'C': 0.01}, {'tol': 0.0001, 'gamma': 100000.0, 'class_weight': 'balanced', 'C': 1e-06}, {'tol': 0.01, 'gamma': 1e-05, 'class_weight': None, 'C': 10000.0}, {'tol': 0.0001, 'gamma': 0.0001, 'class_weight': 'balanced', 'C': 1000000.0}, {'tol': 0.1, 'gamma': 100.0, 'class_weight': None, 'C': 1000000.0}, {'tol': 0.01, 'gamma': 0.0001, 'class_weight': None, 'C': 10.0}, {'tol': 0.0001, 'gamma': 1e-05, 'class_weight': 'balanced', 'C': 1000.0}, {'tol': 0.001, 'gamma': 1e-07, 'class_weight': 'balanced', 'C': 1000.0}], 'split0_test_score': array([0.10016694, 0.10016694, 0.0984975 , 0.14524207, 0.10016694,
       0.0984975 , 0.29215359, 0.10016694, 0.94156928, 0.10016694,
       0.19866444, 0.10016694, 0.66277129, 0.29215359, 0.10016694,
       0.10016694, 0.94991653, 0.97161937, 0.29215359, 0.19866444,
       0.19699499, 0.0984975 , 0.10016694, 0.0984975 , 0.10016694,
       0.19699499, 0.10016694, 0.10016694, 0.94156928, 0.95158598,
       0.10016694, 0.93656093, 0.10016694, 0.19866444, 0.20200334,
       0.0984975 , 0.10016694, 0.19866444, 0.10016694, 0.10016694,
       0.14524207, 0.10016694, 0.10016694, 0.0984975 , 0.94156928,
       0.95158598, 0.10016694, 0.94490818, 0.93989983, 0.94156928]), 'split1_test_score': array([0.10183639, 0.10183639, 0.10016694, 0.0984975 , 0.10183639,
       0.0984975 , 0.10183639, 0.10183639, 0.95492487, 0.10183639,
       0.19866444, 0.10183639, 0.68948247, 0.10183639, 0.10183639,
       0.10183639, 0.96327212, 0.98163606, 0.10183639, 0.0984975 ,
       0.0984975 , 0.0984975 , 0.10183639, 0.0984975 , 0.10183639,
       0.19866444, 0.10183639, 0.10183639, 0.95325543, 0.96327212,
       0.10183639, 0.96661102, 0.10183639, 0.0984975 , 0.10183639,
       0.0984975 , 0.10183639, 0.0984975 , 0.10183639, 0.10183639,
       0.12687813, 0.10183639, 0.10183639, 0.0984975 , 0.95659432,
       0.96327212, 0.10183639, 0.97161937, 0.95826377, 0.95158598]), 'split2_test_score': array([0.10183639, 0.10183639, 0.10016694, 0.0984975 , 0.10183639,
       0.0984975 , 0.10183639, 0.10183639, 0.92654424, 0.10183639,
       0.19866444, 0.10183639, 0.74457429, 0.10183639, 0.10183639,
       0.10183639, 0.93155259, 0.97495826, 0.10183639, 0.0984975 ,
       0.0984975 , 0.0984975 , 0.10183639, 0.0984975 , 0.10183639,
       0.19866444, 0.10183639, 0.10183639, 0.92654424, 0.95158598,
       0.10183639, 0.94156928, 0.10183639, 0.0984975 , 0.10183639,
       0.0984975 , 0.10183639, 0.0984975 , 0.10183639, 0.10183639,
       0.13188648, 0.10183639, 0.10183639, 0.0984975 , 0.93656093,
       0.95158598, 0.10183639, 0.95325543, 0.93989983, 0.92654424]), 'mean_test_score': array([0.10127991, 0.10127991, 0.09961046, 0.11407902, 0.10127991,
       0.0984975 , 0.16527546, 0.10127991, 0.9410128 , 0.10127991,
       0.19866444, 0.10127991, 0.69894268, 0.16527546, 0.10127991,
       0.10127991, 0.94824708, 0.97607123, 0.16527546, 0.13188648,
       0.13132999, 0.0984975 , 0.10127991, 0.0984975 , 0.10127991,
       0.19810796, 0.10127991, 0.10127991, 0.94045632, 0.95548136,
       0.10127991, 0.94824708, 0.10127991, 0.13188648, 0.13522538,
       0.0984975 , 0.10127991, 0.13188648, 0.10127991, 0.10127991,
       0.13466889, 0.10127991, 0.10127991, 0.0984975 , 0.94490818,
       0.95548136, 0.10127991, 0.95659432, 0.94602115, 0.93989983]), 'std_test_score': array([0.00078699, 0.00078699, 0.00078699, 0.0220356 , 0.00078699,
       0.        , 0.08971639, 0.00078699, 0.01159303, 0.00078699,
       0.        , 0.00078699, 0.03405931, 0.08971639, 0.00078699,
       0.00078699, 0.01300314, 0.00416434, 0.08971639, 0.04721915,
       0.04643216, 0.        , 0.00078699, 0.        , 0.00078699,
       0.00078699, 0.00078699, 0.00078699, 0.01093316, 0.0055089 ,
       0.00078699, 0.01314526, 0.00078699, 0.04721915, 0.04721915,
       0.        , 0.00078699, 0.04721915, 0.00078699, 0.00078699,
       0.00775091, 0.00078699, 0.00078699, 0.        , 0.00851255,
       0.0055089 , 0.00078699, 0.01115745, 0.00865684, 0.01029118]), 'rank_test_score': array([25, 25, 45, 24, 25, 46, 15, 25,  9, 25, 13, 25, 12, 15, 25, 25,  5,
        1, 15, 20, 23, 46, 25, 46, 25, 14, 25, 25, 10,  3, 25,  5, 25, 20,
       18, 46, 25, 20, 25, 25, 19, 25, 25, 46,  8,  3, 25,  2,  7, 11],
      dtype=int32)}

Best Params
{'tol': 0.01, 'gamma': 0.001, 'class_weight': None, 'C': 1000.0}
Sleeping 2 minutes to keep job alive if you want to interact!
```

</details>

If you want to debug something, you can set interactive: true to run in interactive mode, and then shell into the pod, connect to the broker:

```bash
$ kubectl exec -it -n flux-operator flux-sample-0-jlsp6 bash
$ sudo -u fluxuser -E $(env) -E HOME=/home/fluxuser flux proxy local:///run/flux/local bash
```

For those interested, when I did this above and looked at the flux jobs submit `flux jobs -a` (to distribute the training)
I saw that indeed, it was distributed across the workers. Notice that we are hitting all the nodes.

```bash
$ flux jobs -a
```
```console
       JOBID USER     NAME       ST NTASKS NNODES     TIME INFO
    ƒAHuVwM1 fluxuser flux-job-+  S      1      -        - 
    ƒ9QwJdRR fluxuser flux-job-+  R      1      1   8.019s flux-sample-3
    ƒ7ebCD4j fluxuser flux-job-+  R      1      1   12.02s flux-sample-0
    ƒ5T6FyQj fluxuser flux-job-+  R      1      1   17.01s flux-sample-1
    ƒ3gfhbD1 fluxuser flux-job-+  R      1      1   21.02s flux-sample-2
    ƒ26RCHAX fluxuser flux-job-+ CD      1      1   5.370s flux-sample-3
```

### Notes

#### Waiting for workers

What I think I saw happening is that a first job was submit that acted as a "nanny" and this
nanny we can ask to wait for workers (see the script). Now, these workers aren't actually the nodes, but 
rather processes on one node (because it's launched by one job). So for that parameter I provided the number of cores.
If you provide the number of workers (and it's more than the number of cores) you will time out!

#### Submitting to Flux

I did this `flux jobs -a` ultiamtely as a sanity check that dask was actually submitting jobs to Flux, and it was.
It does seem like dask doesn't really predict how many workers it will need - it just launches them and the entire
training finishes at some point, and the others can just hang out waiting for more work. I'm sure
this can be tweaked. But seriously, that's so cool! 

#### Cleaning Up

I'm not sure about the logic of cleaning up, because when I used a branch that cleans up when a worker finishes,
I had failed jobs that said they couldn't find the script they were supposed to run. For the time being (until
we understand how to properly cleanup) I have commented out this line. It could be that
the cancel command is not properly going through, or some function call on my part that we are done
needs to happen.

#### Pre-built image

Finally, if you build your image with dask already installed you'll save a bit of time! I also
would play around with the waiting time for the cluster - it likely can be shortened.

### Cleanup

When you are done, clean up:

```bash
$ kubectl delete -f minicluster.yaml
```

Make sure to clean up your temporary tmpdir, because something cached from a previous run could potentially interfere with
a new one! I found that when I didn't clean up, the subsequent run timed out looking for the workers (likely because of some cache?)

```bash
$ sudo rm -rf ./tmp/*
```