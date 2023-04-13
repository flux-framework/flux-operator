#!/usr/bin/env python3
#
#-  step_3_diabetes_postprocessing.py ~~
#

import os
import numpy as np
import pandas as pd

# Obtain demo directory path from environment variable.
demo_dir = os.getenv("DEMO_DIR")

all_coeffs = np.load(os.path.join(demo_dir, "all_coeffs.npy"))

np.set_printoptions(precision=3)
print(all_coeffs)

# Put results into dict with attribute as key.
results = dict()
results["age"] = [all_coeffs[0]]
results["sex"] = [all_coeffs[1]]
results["body_mass_index"] = [all_coeffs[2]]
results["blood_pressure"] = [all_coeffs[3]]
results["total_cholesterol"] = [all_coeffs[4]]
results["ldl_cholesterol"] = [all_coeffs[5]]
results["hdl_cholesterol"] = [all_coeffs[6]]
results["total/hdl_cholesterol"] = [all_coeffs[7]]
results["log_of_serum_triglycerides"] = [all_coeffs[8]]
results["blood_sugar_level"] = [all_coeffs[9]]

df = pd.DataFrame.from_dict(results)

# Print to terminal.
print("Pearson correlation coefficients for each attribute")
print(df.transpose().head(10))

#-  vim:set syntax=python:
