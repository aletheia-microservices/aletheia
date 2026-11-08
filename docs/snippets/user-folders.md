<div class="grid cards compact users" markdown>

-   :material-source-repository: __`blueprint/`__

    Blueprint and its example apps (git submodule). Your own apps go in `blueprint/examples/{app}/`.

-   :material-clipboard-list-outline: __`registry/`__

    Registered applications to be analyzed.

-   :material-tune-variant: __`config/`__

    Per-app detection configs that [suppress warnings](../usage/suppressing-warnings.md).

-   :material-file-document-edit-outline: __`input-models/`__

    [Abstract models](../app-registry/abstract-model.md#guidelines-for-writing-abstract-models) of apps without Blueprint, one folder per app.

-   :material-folder-open-outline: __`output/`__

    Analysis results of each app.

-   :material-script-text-outline: __`scripts/`__

    `app-diagram/` and `schema-diagram/` draw an app's services and schema, `gen-registry/` generates the registry, and `verify/` compares warning counts against a baseline.

-   :material-hammer-wrench: __`Makefile`__

    Build, test, verify, registry, and docs targets (`make help`).

</div>
