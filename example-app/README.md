# example-app

create and deploy an example fly.io app 
```
$ cat adc.json
{
  "universe_domain": "googleapis.com",
  "type": "external_account",
  "audience": "//iam.googleapis.com/projects/1234567890/locations/global/workloadIdentityPools/oidc/providers/fly-io",
  "subject_token_type": "urn:ietf:params:oauth:token-type:jwt",
  "token_url": "https://sts.googleapis.com/v1/token",
  "credential_source": {
    "executable": {
      "command": "fly-auth-plugin"
    }
  }
}

$ fly launch
$ fly secrets set PROJECT_ID=example-project
$ fly secrets set GOOGLE_APPLICATION_CREDENTIALS_JSON="$(cat adc.json)"
```

```
$ curl https://example.fly.dev
Hello from project example-project (projects/1234567890).
```

