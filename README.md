# fly-google-auth

A binary that implements the [Google Application Default Credentials](https://google.aip.dev/auth/4117)
executable provider interface for [Fly.io](https://fly.io/docs/security/openid-connect).

Fly.io provides a local OIDC endpoint at unix:///fly/.api` which returns a JWT
with the following structure:
```
{
  "app_id": "11111111",
  "app_name": "example-app",
  "aud": "https://fly.io/example-org",
  "exp": 1712099653,
  "iat": 1712099053,
  "image": "docker-hub-mirror.fly.io/you/image:latest",
  "image_digest": "sha256:2c1cdaded1b3820020c9dc9fdd1d6e798d6f6ca36861bb6ae64019fad6be9ee3",
  "iss": "https://oidc.fly.io/example-org",
  "jti": "93ca09e1-70e0-477b-a260-1d8fcd4ef4f4",
  "machine_id": "148e21ea7e46e8",
  "machine_name": "example-machine",
  "machine_version": "01HTGGC1TZ2JHK83J4AC0R3VET",
  "nbf": 1712099053,
  "org_id": "11111111",
  "org_name": "example-org",
  "region": "sea",
  "sub": "example-org:example-app:example-machine"
}
```

To use with Google Cloud, configure a Workload Identity pool and provider.
Consider mapping the Fly `sub` field to `google.subject` for use in IAM
bindings.
```
$ gcloud iam workload-identity-pools create $POOL --location=global
$ gcloud iam workload-identity-pools providers create-oidc fly-io \
  --location=global --workload-identity-pool=$POOL \
  --issuer-uri=https://oidc.fly.io/$ORG \
  --attribute-mapping="google.subject=assertion.sub"
```

Then add an IAM binding for your app:
```
$ gcloud projects add-iam-policy-binding $PROJECT \
  --role=roles/browser
  --member=principal://iam.googleapis.com/projects/$NUMBER/locations/global/workloadIdentityPools/$POOL/subject/$ORG:$APP:$MACHINE
```

Finally, configure your Fly.io application with a JSON file like the following:
```
$ export GOOGLE_APPLICATION_CREDENTIALS=adc.json
$ cat adc.json
{
  "universe_domain": "googleapis.com",
  "type": "external_account",
  "audience": "//iam.googleapis.com/projects/$NUMBER/locations/global/workloadIdentityPools/$POOL/providers/$PROVIDER",
  "subject_token_type": "urn:ietf:params:oauth:token-type:jwt",
  "token_url": "https://sts.googleapis.com/v1/token",
  "credential_source": {
    "executable": {
      "command": "fly-google-auth"
    }
  }
}
```

Note: you must set `GOOGLE_EXTERNAL_ACCOUNT_ALLOW_EXECUTABLES=1` to support the
executable provider. Do NOT set this environment variable with a
user-controlled credentials JSON.

