# kube-exec

`kube-exec` is a secure and auditable CLI tool designed to enable users to interact with Kubernetes containers through an interactive shell (`exec`) while maintaining detailed audit logs of all actions.

## Current Features

- Interactive shell into Kubernetes pods (`kubectl exec` functionality).
- Full audit logging in JSON format using `zerolog`.
    - Logs `stdin`, `stdout`, and `stderr` streams with contextual metadata.
    - Metadata includes:
        - `session_id`: Unique identifier for each session.
        - `user`: Executing user.
        - `namespace`: Kubernetes namespace.
        - `pod`: Kubernetes pod name.
        - `container`: Kubernetes container name.
        - `command`: Executed command.
        - Timestamps for each action.
- Lightweight and modular implementation with focus on security.

## TODO: User Authentication

### Functionalities to Implement

1. **LDAP Authentication**:
    - Authenticate users based on LDAP credentials.
    - Map LDAP users to Kubernetes namespaces and roles.

2. **Namespace Access Control**:
    - Validate if the user is authorized to access a specific namespace.
    - Use a `NamespaceExecPolicy` Custom Resource (CR) for defining user permissions.

3. **Session Initiation Authentication**:
    - Validate the user before starting an interactive session.
    - Ensure only authorized users can execute commands in specific containers.

4. **Token-Based Authentication**:
    - Optionally allow authentication using temporary tokens (e.g., JWT or OIDC tokens).
    - Integrate with Kubernetes RBAC.

5. **Audit Enhancements**:
    - Log authentication events (successful and failed attempts).
    - Include the IP address of the user initiating the session.

6. **Pluggable Authentication Backends**:
    - Support additional backends like:
        - Active Directory
        - Custom APIs for authentication

## How to Run

1. Build the CLI:
   ```bash
   go build -o kube-exec main.go
2. Run it !
```shell
./kube-exec exec -n <namespace> -p <pod-name> --interactive -- /bin/bash
```
3. View Logs:

Check logs in `/tmp/kube-exec/interactive.log` (default path).
Each session and its actions are logged in JSON format.