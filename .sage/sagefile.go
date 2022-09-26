package main

import (
	"context"

	"go.einride.tech/sage/sg"
	"go.einride.tech/sage/sgtool"
	"go.einride.tech/sage/tools/sgconvco"
	"go.einride.tech/sage/tools/sggit"
	"go.einride.tech/sage/tools/sggo"
	"go.einride.tech/sage/tools/sggolangcilint"
	"go.einride.tech/sage/tools/sggoreview"
	"go.einride.tech/sage/tools/sgmarkdownfmt"
	"go.einride.tech/sage/tools/sgyamlfmt"
)

func main() {
	sg.GenerateMakefiles(
		sg.Makefile{
			Path:          sg.FromGitRoot("Makefile"),
			DefaultTarget: All,
		},
	)
}

func All(ctx context.Context) error {
	sg.Deps(ctx, ConvcoCheck, FormatMarkdown, FormatYaml, GoGenerate)
	sg.Deps(ctx, GoLint, GoReview, GoTest)
	sg.SerialDeps(ctx, GoModTidy, GitVerifyNoDiff)
	return nil
}

func FormatYaml(ctx context.Context) error {
	sg.Logger(ctx).Println("formatting Yaml files...")
	return sgyamlfmt.Run(ctx)
}

func GoModTidy(ctx context.Context) error {
	sg.Logger(ctx).Println("tidying Go module files...")
	return sg.Command(ctx, "go", "mod", "tidy", "-v").Run()
}

func GoTest(ctx context.Context) error {
	sg.Logger(ctx).Println("running Go tests...")
	return sggo.TestCommand(ctx).Run()
}

func GoReview(ctx context.Context) error {
	sg.Logger(ctx).Println("reviewing Go files...")
	return sggoreview.Command(ctx, "-c", "1", "./...").Run()
}

func GoLint(ctx context.Context) error {
	sg.Logger(ctx).Println("linting Go files...")
	return sggolangcilint.Run(ctx)
}

func FormatMarkdown(ctx context.Context) error {
	sg.Logger(ctx).Println("formatting Markdown files...")
	return sgmarkdownfmt.Command(ctx, "-w", ".").Run()
}

func ConvcoCheck(ctx context.Context) error {
	sg.Logger(ctx).Println("checking git commits...")
	return sgconvco.Command(ctx, "check", "origin/master..HEAD").Run()
}

func GitVerifyNoDiff(ctx context.Context) error {
	sg.Logger(ctx).Println("verifying that git has no diff...")
	return sggit.VerifyNoDiff(ctx)
}

func Stringer(ctx context.Context) error {
	sg.Logger(ctx).Println("installing...")
	_, err := sgtool.GoInstall(ctx, "golang.org/x/tools/cmd/stringer", "v0.1.8")
	return err
}

func GoGenerate(ctx context.Context) error {
	sg.Logger(ctx).Println("generating Go code...")
	sg.Deps(ctx, Stringer)
	return sg.Command(ctx, "go", "generate", "./...").Run()
}
