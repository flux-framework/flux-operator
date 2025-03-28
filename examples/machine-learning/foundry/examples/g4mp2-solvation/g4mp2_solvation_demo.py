#!/usr/bin/env python
# coding: utf-8

# <img src="https://raw.githubusercontent.com/MLMI2-CSSI/foundry/main/assets/foundry-black.png" width=450>

# # Foundry Solvation Energy Quickstart for Beginners
# 
# *Original Paper:* https://doi.org/10.1021/acs.jpca.1c01960
# 
# *Dataset:* https://doi.org/10.18126/c5z9-zej7
# 
# 
# 
# This introduction uses Foundry to:
# 
# 
# 1.   Instantiate and authenticate a Foundry client locally or in the cloud
# 2.   Aggregate data from the G4MP2 solvation database
# 3.   Perform basic data exploration
# 
# 

# [![Open In Colab](https://colab.research.google.com/assets/colab-badge.svg)](https://colab.research.google.com/github/MLMI2-CSSI/foundry/blob/main/examples/g4mp2-solvation/g4mp2_solvation_demo.ipynb)

# This notebook is set up to run locally or as a [Google Colaboratory](https://colab.research.google.com/notebooks/intro.ipynb#scrollTo=5fCEDCU_qrC0) notebook, which allows you to run python code in the browser, or as a [Jupyter](https://jupyter.org/) notebook, which runs locally on your machine.
# 
# The code in the next cell will detect your environment to make sure that only cells that match your environment will run.
# 
# # Environment Set Up
# First we'll need to install Foundry as well as a few other packages. If you're using Google Colab, this code block will install these packages into the Colab environment.
# If you are running locally, it will install these modules onto your machine if you do not already have them. We also have a [requirements file](https://github.com/MLMI2-CSSI/foundry/tree/main/examples/bandgap) included with this notebook. You can run `pip install -r requirements.txt` in your terminal to set up your environment locally.


# We need to import a few packages. We'll be using [Matplotlib](https://matplotlib.org/) to make visualizations of our data, [scikit-learn](https://scikit-learn.org/stable/) to create our model, and [pandas](https://pandas.pydata.org/) and [NumPy ](https://numpy.org/)to work with our data.


import pandas as pd
import matplotlib.pyplot as plt
import seaborn as sns

sns.set_context("poster")


# # Instantiate and Authenticate Foundry
# Once the installations are complete, we can import Foundry.

from foundry import Foundry


# We'll also need to instantiate it. To do so, you'll need a [Globus](https://www.globus.org) account. Once you have your account, you can instantiate Foundry using the code below. When you instantiate Foundry locally, be sure to have your Globus endpoint turned on (you can do that with [Globus Connect Personal](https://www.globus.org/globus-connect-personal)). When you instantiate Foundry on Google Colab, you'll be given a link in the cell's output and asked to enter the provided auth code.


f = Foundry(no_local_server=True, no_browser=True, index="mdf")


# Load the Zeolite Database
# Now that we've installed and imported everything we'll need, it's time to load the data. We'll be loading 1 dataset from Foundry using `f.load` to load the data and then `f.load_data` to load the data into the client.


f.load("10.18126/jos5-wj65", globus=False)
res = f.load_data()

X,y = res['train']
df = pd.concat([X,y], axis=1) # sometimes easier to work with the two together


X.head()


y.head()


# # Data Exploration


sns.set_context('poster')
fig, ax = plt.subplots(figsize=(10,10))

ax.scatter(
    X['u0_atom'],
    y['g4mp2_atom'],
    c=y['sol_acn'],
    s=30,
    alpha=0.5
)

plt.xlim(-1.75, -1.5)
plt.ylim(-1.75, -1.5)

ax.set_xlabel("B3LYP atomization energy at 0K (Ha)")
ax.set_ylabel("G4MP2 atomization energy at 0K (Ha)")
sns.despine()


sns.set_context('poster')
fig, ax = plt.subplots(figsize=(10,10))

ax.scatter(
    y['sol_water'],
    y['sol_acn'],
    c=y['sol_ethanol'],
    s=35,
    alpha=0.3
)

ax.set_xlabel("Solvation Energy in Water (kcal/mol)")
ax.set_ylabel("Solvation Energy in Acetonitrile (kcal/mol)")
sns.despine()

sns.set_context('paper')
sns.pairplot(df[['sol_water', 'sol_acn','sol_ethanol','sol_dmso']])



