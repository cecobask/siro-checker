# siro-checker

## Background

[SIRO](https://siro.ie) is a joint venture between [Vodafone](https://n.vodafone.ie) and [ESB](https://esb.ie). The aim
of SIRO is to roll out high-speed, fibre-to-the-home broadband network across the Republic of Ireland. While most
properties in the [selected areas](https://siro.ie/roll-out) have access to the SIRO network, others don't have access
to such benefits yet. Newly built properties in Ireland frequently fall within the latter group. Many services across
the country heavily rely on the [Eircode](https://www.eircode.ie) system to uniquely identify addresses. If a property
doesn't have an Eircode it can be extremely difficult to avail of basic services. SIRO and most internet service
providers are no exception. Having an Eircode assigned to a property is the first step towards getting access to
fibre-to-the-home broadband network in an area. After that, the next logical step is
to [check if SIRO is available](https://siro.ie/search-your-address). At this point, some people find out that SIRO is
not yet available at their address...

## Motivation

Currently, there is no way to determine what the exact SIRO rollout plans are for individual properties, so the only
solution is to keep checking the availability on the SIRO website. Taking into consideration that many individuals work
from home and/or have families that require decent internet connection at home, the waiting time can seem like an
eternity. The reason I created this web scraper (bot) is to automate the process of continuously checking for
availability related to one's Eircode.

## Usage

1. [Fork the repository](https://github.com/cecobask/siro-checker/fork) to your GitHub account
2. Create a repository secret with name `EIRCODE` and value set to the target Eircode: `Settings` > `Secrets and variables` > `Actions` > `New repository secret`
3. Enable the `check` workflow: `Actions` > `check` > `Enable workflow`
4. Run the `check` workflow manually and verify it worked: `Actions` > `check` > `Run workflow`
5. From now on, GitHub Actions will automatically trigger the `check` workflow every 6 hours
6. Email notification will be sent in the following cases:
    - SIRO is available at your address
    - The `EIRCODE` repository secret is not set
    - Unexpected error occurred ([open new issues](https://github.com/cecobask/siro-checker/issues/new/choose) to report failures)