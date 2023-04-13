#!/usr/bin/env python3
#
#-  submit-workflow.py ~~
#
#-  This gets run by the broker pod first to send tasks to the MongoDB
#       $ python3 submit-workflow.py
#

import os
from fireworks import Firework, LaunchPad, ScriptTask, Workflow

# Set up and reset the LaunchPad using MongoDB URI string.
launchpad = LaunchPad(host = os.getenv("MONGODB_URI"), uri_mode = True)
launchpad.reset("", require_password=False)

# Specify the demo directory via environment variable explicitly. Even though
# this program is running on a login node from the demo directory, the 
# individual Python scripts will execute elsewhere.
demo_dir = os.getenv("DEMO_DIR")

# Create the individual FireWorks and Workflow. The `flux run` will run on a
# batch node, and the Python scripts will execute on compute nodes.
fw1 = Firework(ScriptTask.from_str("flux run -n 1 -c 1 python3 " +
        os.path.join(demo_dir, "step_1_diabetes_preprocessing.py")),
            name = "Step-1")
fw2 = Firework(ScriptTask.from_str("flux run -n 10 -c 1 python3 " +
        os.path.join(demo_dir, "step_2_diabetes_correlation.py")),
            name = "Step-2")
fw3 = Firework(ScriptTask.from_str("flux run -n 1 -c 1 python3 " +
        os.path.join(demo_dir, "step_3_diabetes_postprocessing.py")),
            name = "Step-3")
wf = Workflow([fw1, fw2, fw3], {fw1: fw2, fw2: fw3}, name = "FireWorks demo")

# Store workflow
launchpad.add_wf(wf)
print("Workflow submitted.")

#-  vim:set syntax=python: