#!/usr/bin/env python
# coding: utf-8

# <img src="https://raw.githubusercontent.com/MLMI2-CSSI/foundry/main/assets/foundry-black.png" width=450>

# # Foundry Quantum Monte Carlo ML Quickstart
# 
# *Original Paper:* https://arxiv.org/pdf/2210.06430.pdf
# 
# *Dataset:* https://doi.org/10.18126/wg30-95z0
# 

# This notebook is set up to run locally or as a [Google Colaboratory](https://colab.research.google.com/notebooks/intro.ipynb#scrollTo=5fCEDCU_qrC0) notebook, which allows you to run python code in the browser, or as a [Jupyter](https://jupyter.org/) notebook, which runs locally on your machine.
# 
# The code in the next cell will detect your environment to make sure that only cells that match your environment will run.
# 


no_local_server = True
no_browser = True
globus=False


# # Environment Set Up
# First we'll need to install Foundry as well as a few other packages. If you're using Google Colab, this code block will install these packages into the Colab environment.
# If you are running locally, it will install these modules onto your machine if you do not already have them. We also have a [requirements file](https://github.com/MLMI2-CSSI/foundry/tree/main/examples/bandgap) included with this notebook. You can run `pip install -r requirements.txt` in your terminal to set up your environment locally.
# We need to import a few packages. We'll be using [Matplotlib](https://matplotlib.org/) to make visualizations of our data, [scikit-learn](https://scikit-learn.org/stable/) to create our model, and [pandas](https://pandas.pydata.org/) and [NumPy ](https://numpy.org/)to work with our data.


import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns
import pymatgen as mg
from pymatgen.core import Molecule
import json

sns.set_context("poster")

# # Instantiate and Authenticate Foundry
# Once the installations are complete, we can import Foundry.

from foundry import Foundry


# We'll also need to instantiate it. To do so, you'll need a [Globus](https://www.globus.org) account. Once you have your account, you can instantiate Foundry using the code below. When you instantiate Foundry locally, be sure to have your Globus endpoint turned on (you can do that with [Globus Connect Personal](https://www.globus.org/globus-connect-personal)). When you instantiate Foundry on Google Colab, you'll be given a link in the cell's output and asked to enter the provided auth code.

f = Foundry(no_local_server=no_local_server, no_browser=no_browser, index="mdf")


# Load the Zeolite Database
# Now that we've installed and imported everything we'll need, it's time to load the data. We'll be loading 1 dataset from Foundry using `f.load` to load the data and then `f.load_data` to load the data into the client.

f.load("10.18126/wg30-95z0", globus=globus)
res = f.load_data()


X,y = res['train']
df = pd.concat([X,y], axis=1) # sometimes easier to work with the two together


# # Read in Molecules to PyMatgen

df['mols'] = df['pymatgen'].map(lambda x: Molecule.from_str(x, fmt="json"))


df['mols'].iloc[1]



# # Data Exploration

sns.set_context('poster')
fig, ax = plt.subplots(figsize=(7,7))

ax.scatter(
    y['DMC(HF)'],
    y['DMC(HF)_err'],
    s=30,
    alpha=0.1
)

# plt.xlim(-1.75, -1.5)
# plt.ylim(-1.75, -1.5)

ax.set_xlabel("DMC(HF) (Ha)")
ax.set_ylabel("DMC(HF) error (Ha)")
sns.despine()


sns.set_context('poster')
ax = sns.pairplot(y[['PBE','HF','DMC(HF)','DMC(PBE)','DMC(PBE)_err']], hue='PBE')

fig.savefig('result.png')