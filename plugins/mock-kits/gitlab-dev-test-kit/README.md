# GitLab dev/test kit

A GitLab API v4 contour you can build and test against before anyone gives you
access to a real GitLab — and one you can point CI at without rate limits, flaky
network or a shared instance somebody else is using.

## What you get

| Surface | What appears |
|---------|--------------|
| Mocks | 35 GitLab API v4 endpoints: projects, users, branches, commits, merge requests (approvals, conflicts, diffs, discussions), pipelines, jobs, logs, artifacts, issues, notes, container registry |
| Integrations | a GitLab connection preset pointed at the mocked contour |
| Test cases | five checks for the paths integrations actually get wrong |
| Links | issue → merge request, test case → pipeline |
| Notifications | events for a failed pipeline and a blocked merge request |

## The point: the awkward paths, not the happy one

Most GitLab integrations work against a green pipeline and break on everything
else. This kit ships those cases as first-class mocks:

- **Pagination** — `/projects?page=1` returns a `Link: …rel="next"` header and
  `X-Total-Pages: 2`; page 2 returns no next link. A client that ignores the
  header silently processes half your projects.
- **Rate limiting** — a `429` with `Retry-After: 30` and `RateLimit-Remaining: 0`.
- **A conflicted merge request** — `has_conflicts: true`, and merging it answers
  `405`. A release job must stop, not force it.
- **A failed pipeline** — the state a release gate exists to catch.
- **`401` / `403` / `404` / `409` / `422` / `500`** — including a `500` that
  takes 1.5 seconds, so a client's own timeout gets exercised.

## Installing

```bash
mockarty-cli plugin install mockarty.gitlab-dev-test-kit-1.0.0.zip
mockarty-cli plugin enable mockarty.gitlab-dev-test-kit
```

Then instantiate the **GitLab API v4** mock kit into the namespace you want it
in, and point your client at that namespace's mock base URL.

## What it does not do

It does not reproduce GitLab. There is no repository, no runner and no state
that changes because you pushed something: the merge request that is mergeable
stays mergeable, and the one with conflicts keeps its conflicts. That is what
makes it useful for tests — the same request gives the same answer every time.

If you need behaviour that changes in response to your calls, build it on top:
the mocks are ordinary Mockarty mocks and you can add conditions, chains and
stores to them after instantiating.
