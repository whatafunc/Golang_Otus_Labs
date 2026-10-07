To determine if these tests will fail or work correctly on Git when you push code, we need to simulate the CI/CD process by running the workflow locally. 
Here's a step-by-step guide on how you can do this:

Install GitHub Actions Runner:

You need to install and set up a GitHub Actions runner on your local machine or in a virtual environment. This will allow you to run the workflow steps locally.
Configure the Workflow:

Ensure that the .github/workflows/golang-tests.yml file is correctly configured as shown in the previous example.

Set Up Go Environment:

Install Go 1.21 or later on your local machine.
Set up the GOPATH and ensure it's correctly configured.
Run the Workflow Locally:

Use the GitHub Actions runner to simulate running the workflow locally. You can do this by executing the following commands in your terminal:

# Navigate to your project directory
cd /path/to/your/project

# Install actions-runner-controller if you don't have it installed
helm repo add actions-runner-controller https://actions-runner-controller.github.io/actions-runner-controller
helm install runner actions-runner-controller/actions-runner-controller --namespace actions-runner-system --create-namespace

# Deploy the GitHub Actions runner
kubectl apply -f https://raw.githubusercontent.com/actions/runner-controller/main/deploy/bundle.yaml

Yes, you can simulate a large part of GitHub Actions locally — but **the Kubernetes/Actions Runner Controller instructions you found are overkill for this purpose**.

Actions Runner Controller is for hosting actual **self-hosted GitHub runners in Kubernetes**. GitHub describes ARC specifically as a way to run self-hosted runners in a Kubernetes cluster. [GitHub Docs](https://docs.github.com/en/actions/how-tos/manage-runners/use-actions-runner-controller?utm_source=chatgpt.com)

For your case, use **`act`**. It is specifically designed to read `.github/workflows/*.yml` and execute GitHub Actions locally using Docker. [GitHub](https://github.com/nektos/act?utm_source=chatgpt.com)

### On your Mac

Since you already use Homebrew:

```bash
brew install act
```

Then from:

```text
Golang_Otus_Labs/
```

try:

```bash
act -l
```

That lists the jobs `act` sees in your workflows.

Then you can run a particular job that exists in the yaml file, for example:

```bash
act -j lint
```

or:

```bash
act -j tests
```

`act` uses Docker containers to approximate GitHub's `ubuntu-latest` runner environment. [GitHub](https://github.com/nektos/act?utm_source=chatgpt.com)

### There is one complication with your current workflow

Your YAML currently says:

```yaml
on:
  push:
    branches:
      - hw*
```

and derives:

```yaml
BRANCH=${GITHUB_REF#refs/heads/}
```

then does:

```yaml
working-directory: ${{ env.BRANCH }}
```

Your new branch is `go-rust`, so this old workflow architecture doesn't fit the new repository structure anyway.

Before spending time on local simulation, I think we should first change the CI architecture to something like:

```text
.github/workflows/
├── go-tests.yml
└── rust-tests.yml
```

with your `go-rust` branch running both:

```text
              go-rust
                 │
        ┌────────┴────────┐
        ↓                 ↓
     Go CI             Rust CI
        │                 │
    go test           cargo test
    golangci-lint     cargo clippy
                      cargo fmt --check
```

Then you could test them locally before every push:

```bash
act -j go-tests
act -j rust-tests
```

or run the applicable workflow file directly.

One caveat: `act` is a **very useful approximation**, not a guarantee that GitHub will behave identically. Its documentation explicitly notes that its runner images don't contain everything available on GitHub-hosted runners, and Docker containers aren't identical to GitHub's VMs. [GitHub](https://github.com/nektos/act-docs/blob/main/src/usage/runners.md?utm_source=chatgpt.com)

So your development cycle becomes much nicer:

```text
edit code/workflow
    ↓
go test / cargo test
    ↓
act                    ← simulate CI locally
    ↓
git commit
    ↓
git push
    ↓
real GitHub Actions    ← final authority
```

That is what I'd recommend for this Go→Rust project. **No Kubernetes, Helm, `kubectl`, or Actions Runner Controller is necessary just to test your Actions locally.**

[act — run GitHub Actions locally](https://github.com/nektos/act)

And incidentally, your existing Minikube/Kubernetes experience means ARC could be an interesting experiment later—but it's solving a different problem: providing GitHub with your own runner infrastructure rather than simply checking workflows locally. [GitHub Docs](https://docs.github.com/en/actions/how-tos/manage-runners/use-actions-runner-controller?utm_source=chatgpt.com)