export const resourcesList = {
  data: [
    {
      id: 1,
      name: "ansible-tower-cli",
      catalog: { id: 1, type: "Official" },
      type: "task",
      description:
        "Ansible Tower (formerly ‘AWX’) is a web-based solution that makes Ansible even more easy to use for IT teams of all kinds, It provides the tower-cli(Tower-CLI) command line tool that simplifies the tasks of starting jobs, workflow jobs, manage users, projects etc.",
      versions: [
        { id: 1, version: "1.0" },
        { id: 2, version: "2.0" }
      ],
      tags: [
        { id: 2, tag: "cli" },
        { id: 7, tag: "deploy" },
        { id: 9, tag: "notification" }
      ],
      rating: 0,
      last_updated_at: "2020-04-06T12:30:15.44615Z"
    },
    {
      id: 2,
      name: "argocd",
      catalog: { id: 1, type: "Official" },
      type: "task",
      description:
        "This task syncs (deploys) an Argo CD application and waits for it to be healthy. To do so, it requires the address of the Argo CD server and some form of authentication - either a username/password or an authentication token.",
      versions: [
        { id: 3, version: "1.0" },
        { id: 4, version: "2.0" },
        { id: 5, version: "3.0" }
      ],
      tags: [
        { id: 2, tag: "cli" },
        { id: 5, tag: "azure" }
      ],
      rating: 0,
      last_updated_at: "2020-04-06T12:30:15.450896Z"
    },
    {
      id: 3,
      name: "azure-cli",
      catalog: { id: 1, type: "Official" },
      type: "task",
      description:
        "This task performs operations on Microsoft Azure resources using az.",
      versions: [{ id: 6, version: "1.0" }],
      tags: [{ id: 8, tag: "image-build" }],
      rating: 0,
      last_updated_at: "2020-04-06T12:30:15.453277Z"
    }
  ],
  errors: []
};
