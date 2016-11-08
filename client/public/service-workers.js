// Service workers have very ugly API =/

self.addEventListener(
  'install',
  e => {
    e.waitUntil(
      caches
        .open('expenses-mon')
        .then(cache => {
          return cache
            .addAll([
              '/',
              '/main.js',
              '/main.css',
              '/img/fav/192.png',
              '/img/fav/144.png',
              '/img/fav/96.png',
              '/img/fav/32.png',
              '/manifest.json'
            ])
            .then(() =>
              self.skipWaiting()
            );
        }
      )
    )
  }
);

self.addEventListener(
  'activate',
  event => {
    event.waitUntil(self.clients.claim());
  }
);

self.addEventListener(
  'fetch',
  event => {
    event.respondWith(
      caches
        .match(event.request)
        .then(response => {
          return response || fetch(event.request);
        }
      )
    );
  }
);