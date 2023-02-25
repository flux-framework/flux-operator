# The Flux Operator

<img style="width:50%" alt="Coming Soon" src="_static/images/coming-soon.png">


Welcome to the Flux Operator Documentation!

The Flux Operator is a Kubernetes Cluster [Operator](https://kubernetes.io/docs/concepts/extend-kubernetes/operator/) 
that you can install to your cluster to create and control [Flux Framework](https://flux-framework.org/) "Mini Clusters"
to launch jobs to. With the Flux Operator you can:

1. Create an ephemeral Mini Cluster to run one job
2. Create a persistent Mini Cluster to submit jobs to via user interface, command line, or language SDK.
3. View logs and job information via the same interfaces!

The Flux Operator is currently 🚧️ Under Construction! 🚧️
We are working on core functionality along with providing a library of
examples from HPC. This is a *converged computing* project that aims
to unite the worlds and technologies typical of cloud computing and
high performance computing.

To get started, check out the links below!
Would you like to request a feature or contribute?
[Open an issue](https://github.com/flux-framework/flux-operator/issues).

```{toctree}
:caption: Getting Started
:maxdepth: 2
getting_started/index.md
development/index.md
deployment/index.md
```

```{toctree}
:caption: Tutorials
:maxdepth: 3
getting_started/tutorials/index.md
```


```{toctree}
:caption: Deployment
:maxdepth: 2
deployment/index.md
```

```{toctree}
:caption: About
:maxdepth: 2
about/index.md
```

<script>
// This is a small hack to populate empty sidebar with an image!
document.addEventListener('DOMContentLoaded', function () {
    var currentNode = document.querySelector('.md-sidebar__scrollwrap');
    currentNode.outerHTML =
	'<div class="md-sidebar__scrollwrap">' +
		'<img style="width:100%" src="_static/images/flux-operator.png"/>' +
		
	'</div>';
}, false);

</script>
