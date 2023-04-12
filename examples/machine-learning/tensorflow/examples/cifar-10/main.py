from __future__ import print_function
from __future__ import absolute_import
from __future__ import division

import pickle
import os
import tensorflow.compat.v1 as tf

tf.disable_v2_behavior()

import numpy as np
import sys
from tensorflow_flux import tf_config_from_flux

# ** Modify parameters here for your Minicluster **
cluster, my_job_name, my_task_index = tf_config_from_flux(
    ps_number=1, cluster_size=4, job_name="flux-sample", port_number=2222
)

# original epochs was 1000, would take a few hours!
# I adjusted to 10
epochs = 10

cluster_spec = tf.train.ClusterSpec(cluster)
server = tf.distribute.Server(
    server_or_cluster_def=cluster_spec, job_name=my_job_name, task_index=my_task_index
)

if my_job_name == "ps":
    server.join()
    sys.exit(0)

data_dir = "cifar-10-batches-py"
filelist = [
    os.path.join(data_dir, "data_batch_1"),
    os.path.join(data_dir, "data_batch_2"),
    os.path.join(data_dir, "data_batch_3"),
    os.path.join(data_dir, "data_batch_4"),
    os.path.join(data_dir, "data_batch_5"),
]

data, labels = [], []
is_chief = my_task_index == 0


def unpickle(file):
    with open(file, "rb") as fo:
        obj = pickle.load(fo, encoding="bytes")
    return obj


for f in filelist:
    data_elem = unpickle(f)
    data.append(data_elem[b"data"])
    labels.extend(data_elem[b"labels"])
data = np.vstack(d for d in data)
print("data shape: ", data.shape)


def weight_variable(shape):
    with tf.device("/job:ps/task:0"):
        initial = tf.truncated_normal(shape, stddev=0.1)
        v = tf.Variable(initial)
        return v


def bias_variable(shape):
    with tf.device("/job:ps/task:0"):
        initial = tf.constant(0.1, shape=shape)
        v = tf.Variable(initial)
        return v


def conv2d(x, W):
    return tf.nn.conv2d(x, W, strides=[1, 1, 1, 1], padding="VALID")


def max_pool_2x2(x):
    return tf.nn.max_pool(x, ksize=[1, 2, 2, 1], strides=[1, 2, 2, 1], padding="SAME")


with tf.device("/job:worker/task:{}".format(my_task_index)):
    x = tf.placeholder(tf.float32, shape=[None, 3072], name="x")
    y = tf.placeholder(tf.uint8, shape=[None, 1], name="y")

    # FIRST CONVOLUTIONAL LAYER
    y_one_hot = tf.one_hot(indices=y, depth=10)

    ks = 5
    n_filters1 = 16
    W_conv1 = weight_variable([ks, ks, 3, n_filters1])
    b_conv1 = bias_variable([n_filters1])

    reshaped = tf.reshape(x, [-1, 3, 32, 32])
    transposed = tf.transpose(reshaped, [0, 2, 3, 1])
    x_image = (transposed - 128) / 128

    h_conv1 = tf.nn.relu(conv2d(x_image, W_conv1) + b_conv1)
    h_pool1 = max_pool_2x2(h_conv1)

    # SECOND CONVOLUTIONAL LAYER
    n_filters2 = 64
    W_conv2 = weight_variable([ks, ks, n_filters1, n_filters2])
    b_conv2 = bias_variable([n_filters2])

    h_conv2 = tf.nn.relu(conv2d(h_pool1, W_conv2) + b_conv2)
    h_pool2 = max_pool_2x2(h_conv2)

    # FULLY CONNECTED LAYER
    hidden_neurons = 512
    W_fc1 = weight_variable([5 * 5 * n_filters2, hidden_neurons])
    b_fc1 = bias_variable([hidden_neurons])

    h_pool2_flat = tf.reshape(h_pool2, [-1, 5 * 5 * n_filters2])
    h_fc1 = tf.nn.relu(tf.matmul(h_pool2_flat, W_fc1) + b_fc1)

    # DROPOUT
    keep_prob = tf.placeholder(tf.float32)
    h_fc1_drop = tf.nn.dropout(h_fc1, keep_prob)

    # SOFTMAX
    W_fc2 = weight_variable([hidden_neurons, 10])
    b_fc2 = bias_variable([10])

    y_conv = tf.matmul(h_fc1_drop, W_fc2) + b_fc2
    cross_entropy = tf.nn.softmax_cross_entropy_with_logits(
        logits=y_conv, labels=y_one_hot
    )
    loss = tf.reduce_mean(cross_entropy)
    opt = tf.train.AdamOptimizer(1e-3)
    opt = tf.train.SyncReplicasOptimizer(
        opt,
        replicas_to_aggregate=len(cluster["worker"]),
        total_num_replicas=len(cluster["worker"]),
    )
    global_step = bias_variable([])
    train_step = opt.minimize(loss, global_step=global_step)
    sync_replicas_hook = opt.make_session_run_hook(is_chief)

    y_hat = tf.round(tf.argmax(tf.nn.softmax(y_conv), 1))
    y_hat = tf.cast(y_hat, tf.uint8)
    y_hat = tf.reshape(y_hat, [-1, 1])
    correct_prediction = tf.equal(y_hat, y)
    accuracy = tf.reduce_mean(tf.cast(correct_prediction, tf.float32))


def batch_generator(data, labels, batch_size=32):
    x_batch, y_batch = [], []
    for d, l in zip(data, labels):
        x_batch.append(d)
        y_batch.append(l)
        if len(x_batch) == batch_size:
            yield np.vstack(x_batch), np.vstack(y_batch)
            x_batch = []
            y_batch = []


batch_size = 128
step = 0
sess = tf.train.MonitoredTrainingSession(
    master=server.target, is_chief=is_chief, hooks=[sync_replicas_hook]
)

for i in range(epochs):
    bg = batch_generator(data, labels, batch_size)
    for j, (data_batch, label_batch) in enumerate(bg):
        if (j + i) % len(cluster["worker"]) != my_task_index:
            continue
        _, loss_, acc = sess.run(
            [train_step, loss, accuracy],
            feed_dict={x: data_batch, y: label_batch.reshape(-1, 1), keep_prob: 0.5},
        )
        step += 1
        print(step, my_task_index, loss_, acc)
        sys.stdout.flush()
