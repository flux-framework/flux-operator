#!/usr/bin/env python3

# python ./ray_tune.py --mode train --config-yml path_to/configs/s2ef/200k/forcenet/fn_forceonly.yml --run_dir path_to_run_dir

import numpy as np
from sklearn.datasets import load_digits
from sklearn.model_selection import RandomizedSearchCV
from sklearn.svm import SVC

digits = load_digits()
param_space = {
    "C": np.logspace(-6, 6, 30),
    "gamma": np.logspace(-8, 8, 30),
    "tol": np.logspace(-4, -1, 30),
    "class_weight": [None, "balanced"],
}
model = SVC(kernel="rbf")
search = RandomizedSearchCV(model, param_space, cv=5, n_iter=300, verbose=10)

import joblib
from ray.util.joblib import register_ray

register_ray()
with joblib.parallel_backend("ray"):
    search.fit(digits.data, digits.target)
