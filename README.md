# Pulumi GCP Go: Example Cloud SQL with Internal IP & Private Service Connect

This Pulumi Go template provisions a **Google Cloud SQL (PostgreSQL)** instance with an **internal IP** and **connectivity via Private Service Connect (PSC)**. It also configures a **reserved internal IP address** and a **PSC forwarding rule** to securely access the Cloud SQL instance from another network.

It demonstrates how to:

- Use the Pulumi GCP SDK in Go
- Create a secure Cloud SQL instance (PostgreSQL)
- Set up Private Service Connect (PSC) to access Cloud SQL from on-premises 
- Reserve internal IP addresses and manage forwarding rules

It’s a solid starting point for building **secure, VPC Cloud SQL instance** in production environments.

---

## 📦 Providers

- **Google Cloud Platform (GCP)** via the Pulumi GCP SDK for Go  
  [`github.com/pulumi/pulumi-gcp/sdk/v8/go/gcp`](https://github.com/pulumi/pulumi-gcp)
  - Includes `compute`, `sql`, and `secretmanager` modules
- **Random** provider for secure password generation  
  [`github.com/pulumi/pulumi-random/sdk/v4/go/random`](https://github.com/pulumi/pulumi-random)
- **Pulumi Core SDK** for stack management and outputs  
  [`github.com/pulumi/pulumi/sdk/v3/go/pulumi`](https://github.com/pulumi/pulumi)
---

## ☁️ Resources

- **Cloud SQL Instance** (`gcp.sql.DatabaseInstance`)
  - PostgreSQL
  - Internal IP only
  - Private Service Connect enabled
- **Cloud SQL Database** (`gcp.sql.Database`)
- **Cloud SQL User** (`gcp.sql.User`)
- **Internal IP Address** (`gcp.compute.Address`)
  - Reserved for PSC
- **Forwarding Rule** (`gcp.compute.ForwardingRule`)
  - PSC endpoint to route traffic to Cloud SQL
- **Subnetwork Lookup** (`gcp.compute.LookupSubnetwork`)
  - Fetches subnet where the internal IP and forwarding rule are attached

---

## 🔁 Outputs

- **DB Name**: The name of the Cloud SQL instance
- **DB IP**: The reserved internal IP address used for PSC access
- **DB Connection**: The Cloud SQL instance connection name (used by clients or PSC endpoints)
- **DB Password**: The provisioned or generated password for the database user
- **Forwarding Rule**: The creation timestamp of the PSC forwarding rule (used for verification/troubleshooting)

---

## 📌 When to Use This Template

Use this if:

- You need **secure access to Cloud SQL across VPCs**
- You want to **avoid public IPs** and use **Private Service Connect**
- You're automating Cloud SQL setup with **Pulumi Go**
- You’re building infrastructure where **cloud-native security** is a priority

---

## 🧰 Prerequisites

- Go 1.20 or later installed
- A Google Cloud project with billing enabled
- GCP credentials configured for Pulumi:

```bash
gcloud auth application-default login
```

## 🚀 Usage
1. Scaffold your project
```bash
pulumi new gcp-go
```

2. Set your project ID:
```bash
pulumi config set gcp:project conro-sbx
```

3. Set the Google Cloud Platform region(optional)
```bash
pulumi config set gcp:region europe-west1
```

4. Preview and deploy the resources
```bash
pulumi preview
pulumi up
```

## 🗂️ Project Layout
```bash
├── Pulumi.yaml             # Project metadata
├── Pulumi.<stack>.yaml     # Stack-specific configuration
├── go.mod                  # Go module dependencies
└── main.go                 # Pulumi program: provisions Cloud SQL instance 
```

 ## 📚 Getting Help

 - Pulumi Documentation: https://www.pulumi.com/docs/
 - GCP Provider Reference: https://www.pulumi.com/registry/packages/gcp/
 - Community Slack: https://slack.pulumi.com/
 - GitHub Issues: https://github.com/pulumi/pulumi/issues adjust my readme.
