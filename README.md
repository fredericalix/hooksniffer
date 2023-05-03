# API de stockage et récupération de requêtes JSON

Ce programme utilise Go, Echo framework et PostgreSQL pour créer une API simple permettant de stocker et récupérer des requêtes JSON.

## Configuration

1. Installez les dépendances du projet :

   ```
   go get -u github.com/labstack/echo/v4
   go get -u github.com/lib/pq
   ```

2. Assurez-vous que PostgreSQL est installé et en cours d'exécution.

3. Configurez la variable d'environnement `POSTGRESQL_ADDON_URI` avec l'URI de votre base de données PostgreSQL.

4. Configurez la variable d'environnement `PORT` pour déterminer le port d'écoute du serveur. La valeur par défaut est 8182.

## Lancement de l'API

Exécutez le programme en utilisant la commande suivante :

```
go run main.go
```

## Points de terminaison de l'API

L'API expose les points de terminaison suivants :

### POST /requests

Écoute les requêtes HTTP avec un Content-Type "application/json" et stocke les données dans la base de données.

Exemple d'utilisation avec curl :

```
curl -X POST -H "Content-Type: application/json" -d '{"exemple": "contenu"}' http://localhost:8182/requests
```

### GET /requests

Récupère toutes les requêtes stockées dans la base de données et les renvoie sous forme de tableau JSON.

Exemple d'utilisation avec curl :

```
curl -X GET http://localhost:8182/requests
```

### GET /requests/:id

Récupère la requête ayant l'ID spécifié et la renvoie sous forme d'objet JSON.

Exemple d'utilisation avec curl (remplacez `:id` par l'ID de la requête que vous souhaitez récupérer) :

```
curl -X GET http://localhost:8182/requests/:id
```

## License

Ce projet est sous licence MIT.
