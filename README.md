# Kubelist - Kubernetes Namespace Selector

## How to Run the Application

1. Make sure the `kubelist.json` file containing the namespaces is present at the location `~/.kube/kubelist.json`, with the following structure:

```json
{
  "namespaces": [
    {
      "name": "kube-node-lease"
    },
    {
      "name": "kube-system"
    },
    {
      "name": "kube-public"
    },
    {
      "name": "other"
    }
  ]
}
```

2. Place the kubelist binary in your system's binary or application directory, such as /bin or /Applications, and export the path if necessary.

In the terminal, execute the following command:

```bash
kubelist
```

or

```bash
kubelist <direct-namespace>
```

3. The application will display the list of namespaces available in the kubelist.json file, allow you to select one of them interactively, and set the kubectl context to the selected namespace. Make sure that kubectl is installed on your machine so that the context command works correctly.


> Please note that this is just a basic outline and does not include detailed error handling or other advanced features. You can expand and improve it as needed to meet the requirements of your project.