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
broker.info[0]: quorum-full: quorum->run 8.10777s
     nodes: 4
     cores: 2
   timeout: 60
Fitting model...
Fitting 3 folds for each of 50 candidates, totalling 150 fits
CV Results:
{'mean_fit_time': array([0.15570553, 0.51860126, 0.19759989, 0.64971646, 1.11043668,
       0.70452181, 0.45543218, 0.65649104, 0.18509126, 0.83856575,
       0.5343864 , 0.10854268, 0.48223329, 0.58255227, 0.90180922,
       0.58144879, 0.43473721, 0.50423455, 0.56936741, 0.71839229,
       0.57419117, 0.63295786, 0.70854092, 0.44340865, 0.63092494,
       0.52340007, 0.18849103, 0.63563204, 0.51590824, 0.70667044,
       0.21578566, 1.60080806, 1.33359655, 1.10012269, 1.40437373,
       1.2077709 , 1.61136468, 0.96858056, 0.77815445, 1.86761228,
       1.17164087, 0.82484341, 1.44043525, 0.62390828, 0.1474026 ,
       0.68102757, 0.78141077, 1.01250768, 0.17216452, 0.16020226]), 'std_fit_time': array([0.05973971, 0.02064363, 0.11392156, 0.20227573, 0.23869862,
       0.12768428, 0.26120387, 0.13420861, 0.12622454, 0.34138889,
       0.15459305, 0.03619015, 0.24097666, 0.11489993, 0.38121573,
       0.39901444, 0.08708329, 0.09285345, 0.1375759 , 0.30515549,
       0.03026783, 0.38156088, 0.27626739, 0.07518921, 0.20766368,
       0.09118748, 0.04499471, 0.14277432, 0.25293918, 0.16685543,
       0.16787551, 0.46418934, 0.49988848, 0.33567834, 0.22331225,
       0.4063815 , 0.96429528, 0.28826377, 0.07038422, 0.79501514,
       0.22408482, 0.3940187 , 0.58454422, 0.30810265, 0.04627154,
       0.14212024, 0.2327521 , 0.25167742, 0.03180488, 0.05622747]), 'mean_score_time': array([0.12284168, 0.24133762, 0.14054163, 0.34191608, 0.39344056,
       0.15338906, 0.2731572 , 0.17744438, 0.08095241, 0.30493943,
       0.21205497, 0.10514394, 0.1793685 , 0.18684332, 0.17610041,
       0.28196398, 0.16805506, 0.17218431, 0.28246578, 0.33911284,
       0.251134  , 0.2022237 , 0.19816089, 0.17895484, 0.24495657,
       0.20406127, 0.0819521 , 0.23025791, 0.27222069, 0.17183677,
       0.12053823, 0.48007766, 0.49195663, 0.60704327, 0.32456334,
       0.4532942 , 0.38196683, 0.3058788 , 0.27657493, 0.52431941,
       0.54452038, 0.42188247, 0.27936172, 0.18626753, 0.07600037,
       0.31255293, 0.1613245 , 0.10662246, 0.11645667, 0.07987253]), 'std_score_time': array([0.02947438, 0.13599553, 0.07094923, 0.17235629, 0.19668163,
       0.04341179, 0.12056547, 0.04881971, 0.05108642, 0.13522542,
       0.12704275, 0.05661904, 0.04277962, 0.05244728, 0.0627563 ,
       0.15396325, 0.04543695, 0.03009697, 0.05424151, 0.20868148,
       0.16283435, 0.01805788, 0.03947139, 0.01423403, 0.04341597,
       0.06779754, 0.01885673, 0.0365029 , 0.02512524, 0.0089466 ,
       0.05745425, 0.19493876, 0.07011112, 0.06704884, 0.09924098,
       0.22260647, 0.14120678, 0.10279065, 0.09608677, 0.16628384,
       0.26581479, 0.28444677, 0.08030146, 0.01928522, 0.01555274,
       0.01898241, 0.0056812 , 0.03295954, 0.05005373, 0.02138194]), 'param_tol': masked_array(data=[0.0001, 0.1, 0.001, 0.001, 0.001, 0.1, 0.01, 0.001,
                   0.0001, 0.0001, 0.1, 0.0001, 0.001, 0.0001, 0.1,
                   0.0001, 0.01, 0.1, 0.01, 0.001, 0.01, 0.01, 0.001,
                   0.001, 0.01, 0.0001, 0.001, 0.1, 0.001, 0.001, 0.01,
                   0.01, 0.001, 0.01, 0.1, 0.01, 0.01, 0.0001, 0.0001,
                   0.1, 0.1, 0.0001, 0.1, 0.001, 0.001, 0.0001, 0.1, 0.1,
                   0.1, 0.1],
             mask=[False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False],
       fill_value='?',
            dtype=object), 'param_gamma': masked_array(data=[1e-06, 1e-05, 0.0001, 100000.0, 0.1, 1e-06, 1e-05,
                   10000000.0, 1e-06, 0.1, 100.0, 1e-06, 100.0, 10000.0,
                   100.0, 1e-06, 10000000.0, 0.0001, 0.01, 1000000.0, 0.1,
                   1000.0, 1e-06, 1e-08, 10000000.0, 1e-08, 0.001, 1000.0,
                   1e-05, 0.001, 1e-06, 10.0, 1000.0, 10000.0, 100000.0,
                   10000.0, 0.01, 100.0, 10000000.0, 1000000.0, 1.0,
                   1000.0, 100000000.0, 100000.0, 1e-07, 1e-08, 1e-08,
                   1e-08, 0.001, 0.001],
             mask=[False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False],
       fill_value='?',
            dtype=object), 'param_class_weight': masked_array(data=['balanced', None, 'balanced', None, None, None, None,
                   None, None, 'balanced', None, 'balanced', 'balanced',
                   None, None, 'balanced', None, None, None, 'balanced',
                   'balanced', 'balanced', 'balanced', None, 'balanced',
                   None, 'balanced', None, None, 'balanced', 'balanced',
                   None, 'balanced', None, 'balanced', None, 'balanced',
                   None, None, None, 'balanced', None, None, None,
                   'balanced', None, 'balanced', 'balanced', 'balanced',
                   None],
             mask=[False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False],
       fill_value='?',
            dtype=object), 'param_C': masked_array(data=[1000000.0, 0.01, 100000.0, 0.1, 1000.0, 0.1, 1.0,
                   0.001, 100000.0, 1000.0, 1e-05, 100000.0, 0.1,
                   1000000.0, 100.0, 0.0001, 0.1, 0.0001, 1.0, 1000000.0,
                   1000000.0, 1000.0, 0.0001, 1e-05, 1000000.0, 1e-06,
                   1.0, 1.0, 0.01, 0.0001, 100000.0, 1000000.0, 1000.0,
                   1.0, 0.0001, 100.0, 10000.0, 0.0001, 0.001, 1e-05,
                   1000.0, 0.001, 0.1, 0.0001, 100000.0, 10.0, 0.001,
                   1e-05, 10000.0, 1000000.0],
             mask=[False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False, False, False, False, False, False, False,
                   False, False],
       fill_value='?',
            dtype=object), 'params': [{'tol': 0.0001, 'gamma': 1e-06, 'class_weight': 'balanced', 'C': 1000000.0}, {'tol': 0.1, 'gamma': 1e-05, 'class_weight': None, 'C': 0.01}, {'tol': 0.001, 'gamma': 0.0001, 'class_weight': 'balanced', 'C': 100000.0}, {'tol': 0.001, 'gamma': 100000.0, 'class_weight': None, 'C': 0.1}, {'tol': 0.001, 'gamma': 0.1, 'class_weight': None, 'C': 1000.0}, {'tol': 0.1, 'gamma': 1e-06, 'class_weight': None, 'C': 0.1}, {'tol': 0.01, 'gamma': 1e-05, 'class_weight': None, 'C': 1.0}, {'tol': 0.001, 'gamma': 10000000.0, 'class_weight': None, 'C': 0.001}, {'tol': 0.0001, 'gamma': 1e-06, 'class_weight': None, 'C': 100000.0}, {'tol': 0.0001, 'gamma': 0.1, 'class_weight': 'balanced', 'C': 1000.0}, {'tol': 0.1, 'gamma': 100.0, 'class_weight': None, 'C': 1e-05}, {'tol': 0.0001, 'gamma': 1e-06, 'class_weight': 'balanced', 'C': 100000.0}, {'tol': 0.001, 'gamma': 100.0, 'class_weight': 'balanced', 'C': 0.1}, {'tol': 0.0001, 'gamma': 10000.0, 'class_weight': None, 'C': 1000000.0}, {'tol': 0.1, 'gamma': 100.0, 'class_weight': None, 'C': 100.0}, {'tol': 0.0001, 'gamma': 1e-06, 'class_weight': 'balanced', 'C': 0.0001}, {'tol': 0.01, 'gamma': 10000000.0, 'class_weight': None, 'C': 0.1}, {'tol': 0.1, 'gamma': 0.0001, 'class_weight': None, 'C': 0.0001}, {'tol': 0.01, 'gamma': 0.01, 'class_weight': None, 'C': 1.0}, {'tol': 0.001, 'gamma': 1000000.0, 'class_weight': 'balanced', 'C': 1000000.0}, {'tol': 0.01, 'gamma': 0.1, 'class_weight': 'balanced', 'C': 1000000.0}, {'tol': 0.01, 'gamma': 1000.0, 'class_weight': 'balanced', 'C': 1000.0}, {'tol': 0.001, 'gamma': 1e-06, 'class_weight': 'balanced', 'C': 0.0001}, {'tol': 0.001, 'gamma': 1e-08, 'class_weight': None, 'C': 1e-05}, {'tol': 0.01, 'gamma': 10000000.0, 'class_weight': 'balanced', 'C': 1000000.0}, {'tol': 0.0001, 'gamma': 1e-08, 'class_weight': None, 'C': 1e-06}, {'tol': 0.001, 'gamma': 0.001, 'class_weight': 'balanced', 'C': 1.0}, {'tol': 0.1, 'gamma': 1000.0, 'class_weight': None, 'C': 1.0}, {'tol': 0.001, 'gamma': 1e-05, 'class_weight': None, 'C': 0.01}, {'tol': 0.001, 'gamma': 0.001, 'class_weight': 'balanced', 'C': 0.0001}, {'tol': 0.01, 'gamma': 1e-06, 'class_weight': 'balanced', 'C': 100000.0}, {'tol': 0.01, 'gamma': 10.0, 'class_weight': None, 'C': 1000000.0}, {'tol': 0.001, 'gamma': 1000.0, 'class_weight': 'balanced', 'C': 1000.0}, {'tol': 0.01, 'gamma': 10000.0, 'class_weight': None, 'C': 1.0}, {'tol': 0.1, 'gamma': 100000.0, 'class_weight': 'balanced', 'C': 0.0001}, {'tol': 0.01, 'gamma': 10000.0, 'class_weight': None, 'C': 100.0}, {'tol': 0.01, 'gamma': 0.01, 'class_weight': 'balanced', 'C': 10000.0}, {'tol': 0.0001, 'gamma': 100.0, 'class_weight': None, 'C': 0.0001}, {'tol': 0.0001, 'gamma': 10000000.0, 'class_weight': None, 'C': 0.001}, {'tol': 0.1, 'gamma': 1000000.0, 'class_weight': None, 'C': 1e-05}, {'tol': 0.1, 'gamma': 1.0, 'class_weight': 'balanced', 'C': 1000.0}, {'tol': 0.0001, 'gamma': 1000.0, 'class_weight': None, 'C': 0.001}, {'tol': 0.1, 'gamma': 100000000.0, 'class_weight': None, 'C': 0.1}, {'tol': 0.001, 'gamma': 100000.0, 'class_weight': None, 'C': 0.0001}, {'tol': 0.001, 'gamma': 1e-07, 'class_weight': 'balanced', 'C': 100000.0}, {'tol': 0.0001, 'gamma': 1e-08, 'class_weight': None, 'C': 10.0}, {'tol': 0.1, 'gamma': 1e-08, 'class_weight': 'balanced', 'C': 0.001}, {'tol': 0.1, 'gamma': 1e-08, 'class_weight': 'balanced', 'C': 1e-05}, {'tol': 0.1, 'gamma': 0.001, 'class_weight': 'balanced', 'C': 10000.0}, {'tol': 0.1, 'gamma': 0.001, 'class_weight': None, 'C': 1000000.0}], 'split0_test_score': array([0.93656093, 0.29215359, 0.95158598, 0.10016694, 0.10183639,
       0.29215359, 0.89315526, 0.10016694, 0.93656093, 0.10183639,
       0.10016694, 0.93656093, 0.0984975 , 0.10016694, 0.10016694,
       0.19699499, 0.10016694, 0.29549249, 0.65275459, 0.10016694,
       0.10183639, 0.10016694, 0.19699499, 0.29215359, 0.10016694,
       0.29215359, 0.96828047, 0.10016694, 0.29215359, 0.19866444,
       0.93656093, 0.10016694, 0.10016694, 0.10016694, 0.0984975 ,
       0.10016694, 0.66277129, 0.10016694, 0.10016694, 0.10016694,
       0.20200334, 0.10016694, 0.10016694, 0.10016694, 0.93656093,
       0.29215359, 0.19699499, 0.19699499, 0.96994992, 0.96994992]), 'split1_test_score': array([0.95826377, 0.10183639, 0.96327212, 0.10183639, 0.10183639,
       0.10183639, 0.88313856, 0.10183639, 0.95826377, 0.10183639,
       0.10183639, 0.95826377, 0.0984975 , 0.10183639, 0.10183639,
       0.19866444, 0.10183639, 0.10183639, 0.68447412, 0.10183639,
       0.10183639, 0.10183639, 0.19866444, 0.10183639, 0.10183639,
       0.10183639, 0.98163606, 0.10183639, 0.10183639, 0.19866444,
       0.95659432, 0.10183639, 0.10183639, 0.10183639, 0.10016694,
       0.10183639, 0.68948247, 0.10183639, 0.10183639, 0.10183639,
       0.10183639, 0.10183639, 0.10183639, 0.10183639, 0.95659432,
       0.10183639, 0.0984975 , 0.19866444, 0.98163606, 0.98163606]), 'split2_test_score': array([0.93823038, 0.10183639, 0.95158598, 0.10183639, 0.10183639,
       0.10183639, 0.86477462, 0.10183639, 0.93823038, 0.10183639,
       0.10183639, 0.93823038, 0.0984975 , 0.10183639, 0.10183639,
       0.19866444, 0.10183639, 0.10183639, 0.73789649, 0.10183639,
       0.10183639, 0.10183639, 0.19866444, 0.10183639, 0.10183639,
       0.10183639, 0.97328881, 0.10183639, 0.10183639, 0.19866444,
       0.93823038, 0.10183639, 0.10183639, 0.10183639, 0.10016694,
       0.10183639, 0.74457429, 0.10183639, 0.10183639, 0.10183639,
       0.10183639, 0.10183639, 0.10183639, 0.10183639, 0.93823038,
       0.10183639, 0.0984975 , 0.19866444, 0.97662771, 0.97662771]), 'mean_test_score': array([0.9443517 , 0.16527546, 0.95548136, 0.10127991, 0.10183639,
       0.16527546, 0.88035615, 0.10127991, 0.9443517 , 0.10183639,
       0.10127991, 0.9443517 , 0.0984975 , 0.10127991, 0.10127991,
       0.19810796, 0.10127991, 0.16638843, 0.6917084 , 0.10127991,
       0.10183639, 0.10127991, 0.19810796, 0.16527546, 0.10127991,
       0.16527546, 0.97440178, 0.10127991, 0.16527546, 0.19866444,
       0.94379521, 0.10127991, 0.10127991, 0.10127991, 0.09961046,
       0.10127991, 0.69894268, 0.10127991, 0.10127991, 0.10127991,
       0.13522538, 0.10127991, 0.10127991, 0.10127991, 0.94379521,
       0.16527546, 0.13132999, 0.19810796, 0.97607123, 0.97607123]), 'std_test_score': array([0.0098609 , 0.08971639, 0.0055089 , 0.00078699, 0.        ,
       0.08971639, 0.0117522 , 0.00078699, 0.0098609 , 0.        ,
       0.00078699, 0.0098609 , 0.        , 0.00078699, 0.00078699,
       0.00078699, 0.00078699, 0.09129036, 0.03513343, 0.00078699,
       0.        , 0.00078699, 0.00078699, 0.08971639, 0.00078699,
       0.08971639, 0.0055089 , 0.00078699, 0.08971639, 0.        ,
       0.00907596, 0.00078699, 0.00078699, 0.00078699, 0.00078699,
       0.00078699, 0.03405931, 0.00078699, 0.00078699, 0.00078699,
       0.04721915, 0.00078699, 0.00078699, 0.00078699, 0.00907596,
       0.08971639, 0.04643216, 0.00078699, 0.00478705, 0.00478705]), 'rank_test_score': array([ 5, 18,  4, 29, 26, 18, 10, 29,  5, 26, 29,  5, 50, 29, 29, 14, 29,
       17, 12, 29, 26, 29, 14, 18, 29, 18,  3, 29, 18, 13,  8, 29, 29, 29,
       49, 29, 11, 29, 29, 29, 24, 29, 29, 29,  8, 18, 25, 14,  1,  1],
      dtype=int32)}

Best Params
{'tol': 0.1, 'gamma': 0.001, 'class_weight': 'balanced', 'C': 10000.0}
Sleeping for two minutes to keep job alive if you want to interact!
```

</details>

If you want to debug something, you can set interactive: true to run in interactive mode, and then shell into the pod, connect to the broker:

```bash
$ kubectl exec -it -n flux-operator flux-sample-0-jlsp6 bash
$ sudo -u fluxuser -E $(env) -E HOME=/home/fluxuser flux proxy local:///run/flux/local bash
```

For those interested, when I did this above and looked at the flux jobs submit `flux jobs -a` (to distribute the training)
I saw that indeed, it was distributed across the four workers. Notice that we are hitting all the nodes.

```bash
$ flux jobs -a
```
```console
fluxuser@flux-sample-0:/tmp/workflow$ flux jobs -a
       JOBID USER     NAME       ST NTASKS NNODES     TIME INFO
    ƒ7chV8DD fluxuser dask-work+  R      1      1   1.014s flux-sample-0
    ƒ7XHBkU7 fluxuser dask-work+  R      1      1   1.227s flux-sample-1
    ƒ7T7XmXu fluxuser dask-work+  R      1      1   1.390s flux-sample-2
    ƒ7NmVsdH fluxuser dask-work+  R      1      1   1.548s flux-sample-3
```

It will keep scheduling jobs across workers generally until the task is done. At the end,
I notice that some of the workers just continue to run, and I suspect they will keep going
until the main broker quits (when the job ends) and then the user (you!) can delete
the MiniCluster. If you are interested in what the Dask worker log looks like, take a look
at the flux output files in the present directory (where we submit from):

```bash
$ ls *.out
```
```console
flux-ƒ5G6mJET.out  flux-ƒ5WYQJTH.out  flux-ƒ76S9nYw.out  flux-ƒ8Qi4qTu.out
flux-ƒ5KgKZgj.out  flux-ƒ6u2wJs9.out  flux-ƒ7BebG3H.out  flux-ƒ99HzCxP.out
flux-ƒ5Q5KSA3.out  flux-ƒ6yZM7kB.out  flux-ƒ8KiyFWf.out
```

<detailks>

<summary>Example worker output</summary>

```console
2023-04-23 20:44:31,767 - distributed.nanny - INFO -         Start Nanny at: 'tcp://10.244.0.108:42779'
2023-04-23 20:44:32,542 - distributed.worker - INFO -       Start worker at:   tcp://10.244.0.108:45325
2023-04-23 20:44:32,542 - distributed.worker - INFO -          Listening to:   tcp://10.244.0.108:45325
2023-04-23 20:44:32,542 - distributed.worker - INFO -           Worker name:              FluxCluster-7
2023-04-23 20:44:32,542 - distributed.worker - INFO -          dashboard at:         10.244.0.108:33503
2023-04-23 20:44:32,542 - distributed.worker - INFO - Waiting to connect to:   tcp://10.244.0.109:36893
2023-04-23 20:44:32,543 - distributed.worker - INFO - -------------------------------------------------
2023-04-23 20:44:32,543 - distributed.worker - INFO -               Threads:                          2
2023-04-23 20:44:32,543 - distributed.worker - INFO -                Memory:                   1.86 GiB
2023-04-23 20:44:32,543 - distributed.worker - INFO -       Local Directory: /tmp/workflow/tmp/dask-worker-space/worker-8azkycka
2023-04-23 20:44:32,543 - distributed.worker - INFO - -------------------------------------------------
2023-04-23 20:44:32,796 - distributed.worker - INFO -         Registered to:   tcp://10.244.0.109:36893
2023-04-23 20:44:32,796 - distributed.worker - INFO - -------------------------------------------------
2023-04-23 20:44:32,797 - distributed.core - INFO - Starting established connection to tcp://10.244.0.109:36893
(base) vanessa@vanessa-ThinkPad-T490s:~/Desktop/Code/flux/operator/examples/machine-learning/dask/scikit-learn$ cat flux-ƒ5
flux-ƒ5G6mJET.out  flux-ƒ5KgKZgj.out  flux-ƒ5Q5KSA3.out  flux-ƒ5WYQJTH.out
(base) vanessa@vanessa-ThinkPad-T490s:~/Desktop/Code/flux/operator/examples/machine-learning/dask/scikit-learn$ cat flux-ƒ5G6mJET.out 
2023-04-23 20:44:04,573 - distributed.nanny - INFO -         Start Nanny at: 'tcp://10.244.0.106:43619'
2023-04-23 20:44:05,507 - distributed.worker - INFO -       Start worker at:   tcp://10.244.0.106:43461
2023-04-23 20:44:05,507 - distributed.worker - INFO -          Listening to:   tcp://10.244.0.106:43461
2023-04-23 20:44:05,508 - distributed.worker - INFO -           Worker name:              FluxCluster-1
2023-04-23 20:44:05,508 - distributed.worker - INFO -          dashboard at:         10.244.0.106:34977
2023-04-23 20:44:05,508 - distributed.worker - INFO - Waiting to connect to:   tcp://10.244.0.109:36893
2023-04-23 20:44:05,508 - distributed.worker - INFO - -------------------------------------------------
2023-04-23 20:44:05,509 - distributed.worker - INFO -               Threads:                          2
2023-04-23 20:44:05,509 - distributed.worker - INFO -                Memory:                   1.86 GiB
2023-04-23 20:44:05,509 - distributed.worker - INFO -       Local Directory: /tmp/workflow/tmp/dask-worker-space/worker-de9bramy
2023-04-23 20:44:05,509 - distributed.worker - INFO - -------------------------------------------------
2023-04-23 20:44:05,796 - distributed.worker - INFO -         Registered to:   tcp://10.244.0.109:36893
2023-04-23 20:44:05,796 - distributed.worker - INFO - -------------------------------------------------
2023-04-23 20:44:05,796 - distributed.core - INFO - Starting established connection to tcp://10.244.0.109:36893
2023-04-23 20:44:23,353 - distributed.worker - INFO - Stopping worker at tcp://10.244.0.106:43461. Reason: scheduler-remove-worker
2023-04-23 20:44:23,361 - distributed.nanny - INFO - Closing Nanny gracefully at 'tcp://10.244.0.106:43619'. Reason: scheduler-remove-worker
2023-04-23 20:44:23,436 - distributed.core - INFO - Connection to tcp://10.244.0.109:36893 has been closed.
2023-04-23 20:44:23,442 - distributed.nanny - INFO - Worker closed
[CV 2/3; 4/50] START C=0.1, class_weight=None, gamma=100000.0, tol=0.001........
[CV 3/3; 5/50] START C=1000.0, class_weight=None, gamma=0.1, tol=0.001..........
[CV 2/3; 4/50] END C=0.1, class_weight=None, gamma=100000.0, tol=0.001;, score=0.102 total time=   0.6s
[CV 3/3; 2/50] START C=0.01, class_weight=None, gamma=1e-05, tol=0.1............
[CV 3/3; 5/50] END C=1000.0, class_weight=None, gamma=0.1, tol=0.001;, score=0.102 total time=   1.0s
[CV 3/3; 1/50] START C=1000000.0, class_weight=balanced, gamma=1e-06, tol=0.0001
[CV 3/3; 2/50] END C=0.01, class_weight=None, gamma=1e-05, tol=0.1;, score=0.102 total time=   0.6s
[CV 3/3; 7/50] START C=1.0, class_weight=None, gamma=1e-05, tol=0.01............
[CV 3/3; 1/50] END C=1000000.0, class_weight=balanced, gamma=1e-06, tol=0.0001;, score=0.938 total time=   0.4s
[CV 3/3; 8/50] START C=0.001, class_weight=None, gamma=10000000.0, tol=0.001....
[CV 3/3; 7/50] END C=1.0, class_weight=None, gamma=1e-05, tol=0.01;, score=0.865 total time=   0.6s
[CV 2/3; 9/50] START C=100000.0, class_weight=None, gamma=1e-06, tol=0.0001.....
[CV 3/3; 8/50] END C=0.001, class_weight=None, gamma=10000000.0, tol=0.001;, score=0.102 total time=   0.6s
[CV 3/3; 9/50] START C=100000.0, class_weight=None, gamma=1e-06, tol=0.0001.....
[CV 3/3; 9/50] END C=100000.0, class_weight=None, gamma=1e-06, tol=0.0001;, score=0.938 total time=   0.2s
[CV 2/3; 11/50] START C=1e-05, class_weight=None, gamma=100.0, tol=0.1..........
[CV 2/3; 9/50] END C=100000.0, class_weight=None, gamma=1e-06, tol=0.0001;, score=0.958 total time=   0.5s
[CV 2/3; 12/50] START C=100000.0, class_weight=balanced, gamma=1e-06, tol=0.0001
[CV 2/3; 11/50] END C=1e-05, class_weight=None, gamma=100.0, tol=0.1;, score=0.102 total time=   0.6s
[CV 3/3; 14/50] START C=1000000.0, class_weight=None, gamma=10000.0, tol=0.0001.
[CV 2/3; 12/50] END C=100000.0, class_weight=balanced, gamma=1e-06, tol=0.0001;, score=0.958 total time=   0.3s
[CV 2/3; 15/50] START C=100.0, class_weight=None, gamma=100.0, tol=0.1..........
[CV 3/3; 14/50] END C=1000000.0, class_weight=None, gamma=10000.0, tol=0.0001;, score=0.102 total time=   0.6s
[CV 2/3; 13/50] START C=0.1, class_weight=balanced, gamma=100.0, tol=0.001......
[CV 2/3; 13/50] END C=0.1, class_weight=balanced, gamma=100.0, tol=0.001;, score=0.098 total time=   0.4s
[CV 2/3; 17/50] START C=0.1, class_weight=None, gamma=10000000.0, tol=0.01......
[CV 2/3; 15/50] END C=100.0, class_weight=None, gamma=100.0, tol=0.1;, score=0.102 total time=   1.2s
[CV 3/3; 18/50] START C=0.0001, class_weight=None, gamma=0.0001, tol=0.1........
[CV 2/3; 17/50] END C=0.1, class_weight=None, gamma=10000000.0, tol=0.01;, score=0.102 total time=   0.7s
[CV 2/3; 20/50] START C=1000000.0, class_weight=balanced, gamma=1000000.0, tol=0.001
[CV 2/3; 20/50] END C=1000000.0, class_weight=balanced, gamma=1000000.0, tol=0.001;, score=0.102 total time=   0.4s
[CV 2/3; 22/50] START C=1000.0, class_weight=balanced, gamma=1000.0, tol=0.01...
[CV 3/3; 18/50] END C=0.0001, class_weight=None, gamma=0.0001, tol=0.1;, score=0.102 total time=   0.8s
[CV 3/3; 20/50] START C=1000000.0, class_weight=balanced, gamma=1000000.0, tol=0.001
[CV 2/3; 22/50] END C=1000.0, class_weight=balanced, gamma=1000.0, tol=0.01;, score=0.102 total time=   1.4s
[CV 3/3; 20/50] END C=1000000.0, class_weight=balanced, gamma=1000000.0, tol=0.001;, score=0.102 total time=   1.4s
[CV 2/3; 31/50] START C=100000.0, class_weight=balanced, gamma=1e-06, tol=0.01..
[CV 3/3; 31/50] START C=100000.0, class_weight=balanced, gamma=1e-06, tol=0.01..
[CV 3/3; 31/50] END C=100000.0, class_weight=balanced, gamma=1e-06, tol=0.01;, score=0.938 total time=   0.1s
[CV 2/3; 32/50] START C=1000000.0, class_weight=None, gamma=10.0, tol=0.01......
[CV 2/3; 31/50] END C=100000.0, class_weight=balanced, gamma=1e-06, tol=0.01;, score=0.957 total time=   0.2s
[CV 1/3; 33/50] START C=1000.0, class_weight=balanced, gamma=1000.0, tol=0.001..
[CV 1/3; 33/50] END C=1000.0, class_weight=balanced, gamma=1000.0, tol=0.001;, score=0.100 total time=   1.2s
[CV 2/3; 37/50] START C=10000.0, class_weight=balanced, gamma=0.01, tol=0.01....
[CV 2/3; 32/50] END C=1000000.0, class_weight=None, gamma=10.0, tol=0.01;, score=0.102 total time=   2.2s
[CV 3/3; 36/50] START C=100.0, class_weight=None, gamma=10000.0, tol=0.01.......
[CV 2/3; 37/50] END C=10000.0, class_weight=balanced, gamma=0.01, tol=0.01;, score=0.689 total time=   1.4s
[CV 2/3; 40/50] START C=1e-05, class_weight=None, gamma=1000000.0, tol=0.1......
[CV 3/3; 36/50] END C=100.0, class_weight=None, gamma=10000.0, tol=0.01;, score=0.102 total time=   1.1s
[CV 1/3; 43/50] START C=0.1, class_weight=None, gamma=100000000.0, tol=0.1......
[CV 1/3; 43/50] END C=0.1, class_weight=None, gamma=100000000.0, tol=0.1;, score=0.100 total time=   1.4s
[CV 1/3; 48/50] START C=1e-05, class_weight=balanced, gamma=1e-08, tol=0.1......
[CV 2/3; 40/50] END C=1e-05, class_weight=None, gamma=1000000.0, tol=0.1;, score=0.102 total time=   2.9s
[CV 1/3; 50/50] START C=1000000.0, class_weight=None, gamma=0.001, tol=0.1......
[CV 1/3; 50/50] END C=1000000.0, class_weight=None, gamma=0.001, tol=0.1;, score=0.970 total time=   0.3s
[CV 1/3; 48/50] END C=1e-05, class_weight=balanced, gamma=1e-08, tol=0.1;, score=0.197 total time=   1.4s
2023-04-23 20:44:26,257 - distributed.nanny - INFO - Closing Nanny at 'tcp://10.244.0.106:43619'. Reason: nanny-close-gracefully
2023-04-23 20:44:26,259 - distributed.dask_worker - INFO - End worker
```

</details>

### Notes

#### Submitting to Flux

I did this `flux jobs -a` ultiamtely as a sanity check that dask was actually submitting batch jobs to Flux, and it was.
It does seem like dask doesn't really predict how many of these batch jobs it will need - it just launches them and the entire
training finishes at some point, and the others can just hang out waiting for more work. I'm sure
this can be tweaked. But seriously, that's so cool! 

### Cleanup

When you are done, clean up:

```bash
$ kubectl delete -f minicluster.yaml
```

Make sure to clean up your shared tmpdir!

```bash
$ rm *.out
$ sudo rm -rf ./tmp/*
```