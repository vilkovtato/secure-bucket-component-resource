package main

import (
    "context"

    "github.com/pulumi/pulumi-go-provider/infer"
)

func main() {
    prov, err := infer.NewProviderBuilder().
            WithNamespace("your-org-name").
            WithComponents(
                infer.ComponentF(NewSecureBucket),
            ).
            Build()
    if err != nil {
        panic(err)
    }

    _ = prov.Run(context.Background(), "go-components", "v0.0.1")
}