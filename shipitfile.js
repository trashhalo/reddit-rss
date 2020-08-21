module.exports = shipit => {
  require('shipit-deploy')(shipit)

  shipit.initConfig({
    default: {
      workspace: '/tmp/reddit-rss',
      deployTo: '/usr/src/reddit-rss',
      repositoryUrl: 'https://github.com/trashhalo/reddit-rss.git',
      keepWorkspaces: false,
    },
    production: {
      servers: process.env.SERVER,
    },
  });

  shipit.on('fetched', () => {
    shipit.start('app:build');
  });

  shipit.blTask('app:build', async () => {
    await shipit.local('go test ./...', {cwd: shipit.workspace});
    await shipit.local('CGO_ENABLED=0 go build -a -ldflags \'-extldflags "-static"\' ./cmd/reddit-rss', {cwd: shipit.workspace});
  });

  shipit.on('published', () => {
    shipit.start('app:restart');
  });

  shipit.blTask('app:restart', async () => {
    await shipit.remote('sudo systemctl restart reddit-rss');
  });
}
