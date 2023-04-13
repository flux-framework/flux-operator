#!/usr/bin/env python3
#
#-  step_1_diabetes_preprocessing.py ~~
#
#   This is the first step of the three step machine learning workflow.
#

import os
import numpy as np
import pandas as pd
from sklearn.datasets import load_diabetes

# Obtain demo directory path from environment variable.
demo_dir = os.getenv("DEMO_DIR")

# Load public diabetes dataset.
x_diabetes, y_diabetes = load_diabetes(return_X_y=True)

# Dataset info on the docs page: https://scikit-learn.org/stable/datasets/toy_dataset.html#diabetes-dataset
# https://www4.stat.ncsu.edu/~boos/var.select/diabetes.html

# Note: Each of these 10 feature variables have been mean centered and scaled by the standard deviation
# times the square root of n_samples (i.e. the sum of squares of each column totals 1).

# Number of Instances:
#     442
# Number of Attributes:
#     First 10 columns are numeric predictive values
# Target:
#     Column 11 is a quantitative measure of disease progression one year after baseline
# Attribute Information:
#         age age in years
#         sex
#         bmi body mass index
#         bp average blood pressure
#         s1 tc, total serum cholesterol
#         s2 ldl, low-density lipoproteins
#         s3 hdl, high-density lipoproteins
#         s4 tch, total cholesterol / HDL
#         s5 ltg, possibly log of serum triglycerides level
#         s6 glu, blood sugar level

# Take a look at the data, but remember it has been mean centered and scaled.
print(pd.DataFrame(x_diabetes).describe())

# Save attributes.
np.save(os.path.join(demo_dir, "x_diabetes.npy"), x_diabetes)
print("wrote x_diabetes file")

# Save disease progression measure.
np.save(os.path.join(demo_dir, "y_diabetes.npy"), y_diabetes)
print("wrote y_diabetes file")

#-  vim:set syntax=python:
