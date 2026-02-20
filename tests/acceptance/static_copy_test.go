package acceptance

import (
	"path/filepath"
	"testing"
)

func TestStaticCopy_CopiesAllFiles(t *testing.T) {
	ctx, projectDir := SetupProject(t, "static", "testproject")
	templateDir := filepath.Join(ProjRoot(), "testdata/share/proj/static")

	VerifyFileHash(t, projectDir, "readme.md", filepath.Join(templateDir, "readme.md"))
	VerifyFileHash(t, projectDir, "readme-example1.md", filepath.Join(templateDir, "src/example1.md"))
	VerifyFileHash(t, projectDir, "src/example2.md", filepath.Join(templateDir, "src/example2.md"))
	VerifyFileHash(t, projectDir, "src/foo/bar/example3.md", filepath.Join(templateDir, "src/example3.md"))
	VerifyFileHash(t, projectDir, "src/example-4-testproject.md", filepath.Join(templateDir, "src/example4.md"))
	VerifyFileHash(t, projectDir, "src/example5.md", filepath.Join(templateDir, "src/example5.md"))
	VerifyFileHash(t, projectDir, "example6.sh", filepath.Join(templateDir, "src/example6.sh"))

	_ = ctx // ctx is used for cleanup via t.Cleanup
}
