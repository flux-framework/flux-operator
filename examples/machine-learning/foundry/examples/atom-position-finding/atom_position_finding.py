#!/usr/bin/env python
# coding: utf-8

# # Installing Foundry
# First we'll need to install Foundry. We'll also be installing [Matplotlib](https://matplotlib.org/) for our visualizations. If you're using Google Colab, this code block will install this package into the Colab environment.
# 
# 
# If you are running locally, it will install this module onto your machine if you do not already have it. We also have a [requirements file](https://github.com/MLMI2-CSSI/foundry/tree/main/examples/atom-position-finding) included with this notebook. You can run `pip install -r requirements.txt` in your terminal to set up your environment locally.


# # Importing Packages
# Now we can import Foundry and Matplotlib so we can import the data and visualize it.

# In[9]:


from foundry import Foundry
import matplotlib.pyplot as plt 

# # Instantiating Foundry
# To instantiate Foundry, you'll need a [Globus](https://www.globus.org) account. Once you have your account, you can instantiate Foundry using the code below. When you instantiate Foundry locally, be sure to have your Globus endpoint turned on (you can do that with [Globus Connect Personal](https://www.globus.org/globus-connect-personal)). When you instantiate Foundry on Google Colab, you'll be given a link in the cell's output and asked to enter the provided auth code.

f = Foundry(index="mdf", no_local_server=True, no_browser=True)

dataset_doi = '10.18126/e73h-3w6n'

# download the data 
f.load(dataset_doi, download=True, globus=False)

# load the HDF5 image data into a local object
res = f.load_data()

# using the 'train' split, 'input' or 'target' type, and Foundry Keys specified by the dataset publisher
# we can grab the atom images, metadata, and coorinates we desire
imgs = res['train']['input']['imgs']
desc = res['train']['input']['metadata']
coords = res['train']['target']['coords']

n_images = 3
offset = 150
key_list = list(res['train']['input']['imgs'].keys())[0+offset:n_images+offset]

fig, axs = plt.subplots(1, n_images, figsize=(20,20))
for i in range(n_images):
    axs[i].imshow(imgs[key_list[i]])
    axs[i].scatter(coords[key_list[i]][:,0], coords[key_list[i]][:,1], s = 20, c = 'r', alpha=0.5)

fig.savefig("result.png")