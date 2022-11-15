# Configuration file for the Sphinx documentation builder.
#
# This file only contains a selection of the most common options. For a full
# list see the documentation:
# http://www.sphinx-doc.org/en/master/config

# -- Path setup --------------------------------------------------------------

# If extensions (or modules to document with autodoc) are in another directory,
# add these directories to sys.path here. If the directory is relative to the
# documentation root, use os.path.abspath to make it absolute, like shown here.
#
import os
import sys
   
sys.path.insert(0, os.path.abspath(os.path.dirname(os.path.abspath("."))))
from distutils.version import LooseVersion  # noqa

import sphinx_material  # noqa
from recommonmark.transform import AutoStructify  # noqa

FORCE_CLASSIC = os.environ.get("SPHINX_MATERIAL_FORCE_CLASSIC", False)
FORCE_CLASSIC = FORCE_CLASSIC in ("1", "true")

# -- Project information -----------------------------------------------------

project = "Flux Operator"
html_title = "Flux Operator"

copyright = "2022, Flux Developers"
author = "@vsoch"

# The full version, including alpha/beta/rc tags
release = LooseVersion(sphinx_material.__version__).vstring

# -- General configuration ---------------------------------------------------

# Add any Sphinx extension module names here, as strings. They can be
# extensions coming with Sphinx (named 'sphinx.ext.*') or your custom
# ones.
extensions = [
    "myst_parser",
    "sphinx.ext.autosummary",
    "sphinx.ext.autodoc",
    "sphinx.ext.doctest",
    "sphinx.ext.extlinks",
    "sphinx.ext.intersphinx",
    "sphinx.ext.todo",
    "sphinx.ext.mathjax",
    "sphinx.ext.viewcode",
    "nbsphinx",
    "sphinx_markdown_tables",
    "sphinx_copybutton",
    "sphinx_search.extension",
]


autosummary_generate = True
autoclass_content = "class"

# Add any paths that contain templates here, relative to this directory.
templates_path = ["_templates"]

# List of patterns, relative to source directory, that match files and
# directories to ignore when looking for source files.
# This pattern also affects html_static_path and html_extra_path.
exclude_patterns = [
    "_build",
    "Thumbs.db",
    ".DS_Store",
    "env",
    "README.md",
    ".github",
    ".circleci",
]

# -- Options for HTML output -------------------------------------------------

# The theme to use for HTML and HTML Help pages.  See the documentation for
# a list of builtin themes.
#
# Add any paths that contain custom static files (such as style sheets) here,
# relative to this directory. They are copied after the builtin static files,
# so a file named 'default.css' will overwrite the builtin 'default.css'.
html_static_path = ["_static"]

# -- HTML theme settings ------------------------------------------------

html_show_sourcelink = True
html_sidebars = {
    "**": ["logo-text.html", "globaltoc.html", "localtoc.html", "searchbox.html"]
}

# Allows us to add to the default template
templates_path = ["_templates"]

extensions.append("sphinx_material")
html_theme_path = sphinx_material.html_theme_path()
html_context = sphinx_material.get_html_context()
html_theme = "sphinx_material"
html_css_files = ["custom.css"]

# Custom sphinx material variables
theme_logo_icon = "images/the-operator.jpg"

html_theme_options = {
    "base_url": "http://flux-framework.github.io/flux-operator/",
    "repo_url": "https://github.com/flux-framework/flux-operator/",
    "repo_name": "Flux Operator",
    "html_minify": False,
    "html_prettify": True,
    "css_minify": False,
    # https://fonts.google.com/icons?icon.query=cycle
    "logo_icon": "cycle",
    "repo_type": "github",
    "globaltoc_depth": 2,
    # red, pink, purple, deep-purple, indigo, blue, light-blue, cyan, teal, green, light-green, lime, yellow, amber, orange, deep-orange, brown, grey, blue-grey, and white.
    "color_primary": "blue",
    # red, pink, purple, deep-purple, indigo, blue, light-blue, cyan, teal, green, light-green, lime, yellow, amber, orange, and deep-orange.
    "color_accent": "blue",
    "touch_icon": "images/flux-operator.jpg",
    "theme_color": "#036291",
    "master_doc": False,
    "nav_links": [
        {
            "href": "https://flux-framework.org/",
            "internal": False,
            "title": "Flux Framework",
        },
        {
            "href": "https://github.com/flux-framework",
            "internal": False,
            "title": "Flux Framework on GitHub",
        },
        {
            "href": "https://github.com/flux-framework/flux-operator",
            "internal": False,
            "title": "Flux Operator on GitHub",
        },
    ],
    "heroes": {
        "index": "Flux Operator",
        "customization": "Flux Operator",
    },
    # Include the version dropdown top right? (e.g., if we use readthedocs)
    "version_dropdown": False,
    # Format of this is dict with [label,path]
    # Since we are rendering on gh-pages without readthedocs, we don't
    # have versions
    # "version_json": "_static/versions.json",
    # "version_info": {
    #    "Release": "https://online-ml.github.io/viz/",
    #    "Development": "https://online-ml.github.io/viz/devel/",
    #    "Release (rel)": "/viz/",
    #    "Development (rel)": "/viz/devel/",
    # },
    # Do NOT strip these classes from tables!
    "table_classes": ["plain"],
}

if FORCE_CLASSIC:
    print("!!!!!!!!! Forcing classic !!!!!!!!!!!")
    html_theme = "classic"
    html_theme_options = {}
    html_sidebars = {"**": ["globaltoc.html", "localtoc.html", "searchbox.html"]}

language = "en"
html_last_updated_fmt = ""

todo_include_todos = True
html_favicon = "images/favicon.ico"

html_use_index = True
html_domain_indices = True

nbsphinx_execute = "always"
nbsphinx_kernel_name = "python3"

extlinks = {
    "duref": (
        "http://docutils.sourceforge.net/docs/ref/rst/" "restructuredtext.html#%s",
        "",
    ),
    "durole": ("http://docutils.sourceforge.net/docs/ref/rst/" "roles.html#%s", ""),
    "dudir": ("http://docutils.sourceforge.net/docs/ref/rst/" "directives.html#%s", ""),
}


# Enable eval_rst in markdown
def setup(app):
    app.add_config_value(
        "recommonmark_config",
        {"enable_math": True, "enable_inline_math": True, "enable_eval_rst": True},
        True,
    )
    app.add_transform(AutoStructify)
    app.add_object_type(
        "confval",
        "confval",
        objname="configuration value",
        indextemplate="pair: %s; configuration value",
    )
