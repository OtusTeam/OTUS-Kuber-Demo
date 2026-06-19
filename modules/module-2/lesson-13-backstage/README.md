# Backstage with GitLab Integration

This setup provides a Docker Compose configuration for running Backstage with GitLab integration.

## Current Status

The current setup uses a simple HTTP server as a placeholder. This demonstrates that the Docker configuration works correctly.

## Proper Backstage Setup

For a complete Backstage application with GitLab integration, you would need to:

1. Create a Backstage application using `npx @backstage/create-app`
2. Configure GitLab integration in the app
3. Build and deploy the application

## Docker Compose Configuration

The `docker-compose.yml` file is configured to:
- Build from the included Dockerfile
- Map port 7007 for the Backstage UI
- Pass the GitLab token as an environment variable
- Mount the configuration file

## Configuration

The `app-config.yaml` file includes GitLab integration settings that would be used by a full Backstage application.

## Next Steps

To create a full Backstage application:

1. On your local machine (not in the container), run:
   ```bash
   npx @backstage/create-app@latest backstage-app
   ```
   
2. Follow the prompts to create the application

3. Configure GitLab integration in the generated app

4. Update the Dockerfile to build and run the actual Backstage application

## Environment Variables

Set your GitLab token before running:
```bash
export GITLAB_TOKEN=your_actual_gitlab_token
```

Then start the service:
```bash
docker-compose up --build
```


https://github.com/backstage/backstage/tree/master/contrib/docker/cookiecutter-with-jinja2-extensions

Для интеграции шаблонов Cookiecutter в Backstage и их использования применяется встроенное действие (action) fetch:cookiecutter Fetch And Run Cookiecutter - Backstage Scaffolder Action Function createFetchAndRunCookiecutterAction - Backstage. Весь процесс состоит из настройки бэкенда, создания YAML-файла шаблона и его запуска через интерфейс.

Шаг 1. Подготовка бэкенда Backstage
Убедитесь, что модуль Cookiecutter подключен в вашем проекте.Перейдите в директорию бэкенда и добавьте необходимый пакет Builtin actions | Backstage Software Catalog and Developer ...:bashyarn --cwd packages/backend add @backstage/plugin-scaffolder-backend-module-cookiecutter
Используйте код с осторожностью.Откройте файл packages/backend/src/index.ts и зарегистрируйте модуль:

```js
typescriptimport { createRouter } from '@backstage/plugin-scaffolder-backend';
// ...
// Внутри функции createRouter или настройки бэкенда добавьте:
backend.add(import('@backstage/plugin-scaffolder-backend-module-cookiecutter/alpha'));
```
Используйте код с осторожностью.(Примечание: если у вас старая конфигурация, используйте scaffolderEnv.addActions(createFetchCookiecutterAction(...))).


Шаг 2. Создание файла шаблона (template.yaml)Создайте YAML-файл, который будет описывать интерфейс ввода параметров для пользователей и последовательность шагов Шаблоны программного обеспечения Backstage:yaml

```yaml
apiVersion: scaffolder.backstage.io/v1beta3
kind: Template
metadata:
  name: cookiecutter-react-template
  title: Cookiecutter React Microservice
  description: Создание микросервиса на React с использованием Cookiecutter
spec:
  owner: user:guest
  type: service
  
  # Описание параметров, которые будет заполнять пользователь в UI
  parameters:
    - title: Заполните данные проекта
      required:
        - name
        - description
      properties:
        name:
          title: Имя проекта
          type: string
          description: Название создаваемого сервиса
        description:
          title: Описание
          type: string
        version:
          title: Версия
          type: string
          default: '0.1.0'
```

### Шаги выполнения, которые выполнит Backstage
```yaml
  steps:
    - id: fetch
      name: Генерация файлов из Cookiecutter
      action: fetch:cookiecutter
      input:
        url: https://github.com
        values:
          project_name: ${{ parameters.name }}
          description: ${{ parameters.description }}
          version: ${{ parameters.version }}

    - id: publish
      name: Публикация в Git
      action: publish:github # или publish:gitlab
      input:
        allowedHosts: ['github.com']
        description: ${{ parameters.description }}
        repoUrl: ://github.com{{ parameters.name }}

    - id: register
      name: Регистрация в Каталоге Backstage
      action: catalog:register
      input:
        repoContentsUrl: ${{ steps['publish'].output.remoteUrl }}
        catalogInfoPath: '/catalog-info.yaml'
```

Шаг 3. Регистрация шаблона в BackstageЧтобы шаблон появился в портале, его нужно зарегистрировать:Перейдите в раздел Catalog в Backstage.Нажмите кнопку Create (или Register Entity).Укажите ссылку на ваш template.yaml (например, github.com).Шаг 4. Использование шаблонаВ меню Create (Scaffolder) выберите созданный вами шаблон Cookiecutter React Microservice.В появившемся визарде (мастере) заполните поля формы, указав имя проекта и описание.Нажмите Next. Backstage запустит действие fetch:cookiecutter, скачает ваш шаблон, применит к нему указанные параметры, создаст новый репозиторий, запушит код и автоматически добавит его в каталог Backstage.
