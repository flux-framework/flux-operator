# Project Governance

This repository (the `flux-operator`) is a component of the **Flux Framework ecosystem**.

While this project operates under the high-level guidance of the [Flux Framework Project](https://flux-framework.readthedocs.io/projects/flux-core/en/latest/guide/support.html), it maintains its own list of maintainers, release cadence, and technical decision-making processes tailored to the specific language and domain (e.g., Python, Go/Kubernetes, REST APIs) of the component.

## 1. Roles and Responsibilities

### 1.1 Contributors

Contributors are community members who submit patches, file issues, improve documentation, or answer questions on community channels. Anyone can be a contributor. There is no expectation of long-term commitment.

### 1.2 Maintainers
Maintainers are contributors who have shown dedication to the project through consistent, high-quality contributions. Maintainers have **Write Access** to the repository.

**Responsibilities:**
*   Reviewing and merging Pull Requests (PRs).
*   Triaging and managing Issues.
*   Cutting releases and publishing artifacts (e.g., to PyPI, Docker Hub, or GitHub Releases).
*   Ensuring code quality, test coverage, and security standards are met.
*   Participating in technical discussions and design reviews.

**Becoming a Maintainer:**
Maintainership is earned by merit. Existing maintainers may nominate a contributor who has demonstrated:

1.  A solid understanding of the project's codebase and goals.
2.  A history of high-quality pull requests.
3.  A collaborative and respectful attitude toward other community members.
4.  Commitment to the project's [Code of Conduct](CODE_OF_CONDUCT.md).

New maintainers are added by a consensus vote of the existing maintainers.

## 2. Decision Making

### 2.1 Lazy Consensus
To ensure velocity, this project operates primarily on **Lazy Consensus**.
*   Changes are proposed via Pull Requests or Issues.
*   If a Maintainer reviews and approves the change, and no other Maintainer objects within a reasonable timeframe (usually 48 hours for non-trivial changes), the change is accepted.
*   Silence is not interpreted as consent.

### 2.2 Voting
In the rare event that consensus cannot be reached:
1.  Maintainers will attempt to resolve the disagreement through discussion on GitHub or the developer meeting.
2.  If a deadlock persists, a simple majority vote of the current Maintainers of this repository decides the outcome.
3.  In the event of a tie, the issue is escalated to the Flux Framework developers.

## 3. Contribution Process

### 3.1 Developer Certificate of Origin (DCO)

To ensure the project can legally redistribute your contributions, all commits must be signed off. This certifies that you wrote the code or have the right to contribute it.
*   Git command: `git commit -s ...`
*   This adds a `Signed-off-by: Name <email>` trailer to your commit message.

### 3.2 Code Review
*   All code changes must be submitted via Pull Request.
*   **Requirement:** At least **one** approval from a Maintainer is required to merge.

## 4. Code of Conduct

This project adheres to a **Code of Conduct**.
*   [Link to Code of Conduct](https://github.com/flux-framework/flux-operator/blob/master/.github/CODE_OF_CONDUCT.md)
*   Instances of abusive, harassing, or otherwise unacceptable behavior may be reported by contacting the project lead [sochat1@llnl.gov](mailto:sochat1@llnl.gov).

## 5. Current Maintainers

*   **Vanessa Sochat** (@vsoch) - *Lawrence Livermore National Laboratory*
