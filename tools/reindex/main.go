// Command reindex regenerates examples/plugins/registry-index/index.json from
// every example plugin directory: it packs each bundle (the same BundleDir the
// installer and the registry-drift guard use), reads its manifest, and writes a
// registry entry with the freshly-computed sha256. Adding a plugin is then just
// "drop a dir, run reindex" — the guard TestExampleBundlesMatchRegistryIndex
// stays green without hand-editing the index.
//
//	go run ./examples/plugins/tools/reindex
package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"mockarty/internal/plugins"
)

// downloadURL is the convention for the (to-be-created) public registry repo.
func downloadURL(id, version string) string {
	return fmt.Sprintf("https://github.com/mockarty/plugins-registry/releases/download/%s-v%s/%s-%s.zip", id, version, id, version)
}

type entry struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Version     string `json:"version"`
	Description string `json:"description"`
	AuthorName  string `json:"author_name,omitempty"`
	Homepage    string `json:"homepage,omitempty"`
	DownloadURL string `json:"download_url"`
	SHA256      string `json:"sha256"`
}

type index struct {
	Kind    string  `json:"kind"`
	Version int     `json:"version"`
	Plugins []entry `json:"plugins"`
}

func main() {
	root := "examples/plugins"
	if len(os.Args) > 1 {
		root = os.Args[1]
	}
	des, err := os.ReadDir(root)
	must(err)
	idx := index{Kind: "mockarty/plugin-registry", Version: 1}
	for _, de := range des {
		if !de.IsDir() || de.Name() == "registry-index" || de.Name() == "tools" {
			continue
		}
		dir := filepath.Join(root, de.Name())
		if _, err := os.Stat(filepath.Join(dir, plugins.ManifestFileName)); err != nil {
			continue
		}
		raw, err := plugins.BundleDir(dir)
		if err != nil {
			fmt.Fprintf(os.Stderr, "SKIP %s: pack failed (build its .wasm?): %v\n", de.Name(), err)
			os.Exit(1)
		}
		b, err := plugins.ReadBundle(raw)
		must(err)
		m := b.Manifest
		idx.Plugins = append(idx.Plugins, entry{
			ID: m.ID, Name: m.Name, Version: m.Version, Description: m.Description,
			AuthorName: m.Author.Name, Homepage: m.Homepage,
			DownloadURL: downloadURL(m.ID, m.Version), SHA256: b.SHA256,
		})
	}
	sort.Slice(idx.Plugins, func(i, j int) bool { return idx.Plugins[i].ID < idx.Plugins[j].ID })
	out, err := json.MarshalIndent(idx, "", "  ")
	must(err)
	out = append(out, '\n')
	must(os.WriteFile(filepath.Join(root, "registry-index", "index.json"), out, 0o644))
	fmt.Printf("reindexed %d plugins → %s/registry-index/index.json\n", len(idx.Plugins), root)
}

func must(err error) {
	if err != nil {
		panic(err)
	}
}
