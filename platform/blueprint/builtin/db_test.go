package builtin

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/boltdb/bolt"
	"github.com/micromdm/micromdm/platform/blueprint"
	profile "github.com/micromdm/micromdm/platform/profile/builtin"
)

func TestSave(t *testing.T) {
	db := setupDB(t)
	bp := &blueprint.Blueprint{}
	bp.ApplyAt = []string{"Enroll"}

	if err := db.Save(bp); err == nil {
		t.Fatal("blueprints are required to have an UUID and Name")
	}

	bp.UUID = ""
	bp.Name = "blueprint"
	if err := db.Save(bp); err == nil {
		t.Fatal("blueprints are required to have an UUID")
	}

	bp.UUID = "a-b-c-d"
	bp.Name = ""
	if err := db.Save(bp); err == nil {
		t.Fatal("blueprints are required to have a Name")
	}

	bp.UUID = "a-b-c-d"
	bp.Name = "blueprint"
	if err := db.Save(bp); err != nil {
		t.Fatalf("saving blueprint in datastore: %v", err)
	}

	bp.UUID = "e-f-g-h"
	bp.Name = "blueprint"
	if err := db.Save(bp); err == nil {
		t.Fatal("blueprint names must be unique")
	}

	bp.UUID = "e-f-g-h"
	bp.Name = "blueprint2"
	if err := db.Save(bp); err != nil {
		t.Fatalf("saving blueprint2 in datastore: %v", err)
	}

	byName, err := db.BlueprintByName("blueprint")
	if err != nil {
		t.Fatalf("getting blueprint by Name: %v", err)
	}
	if byName == nil || byName.UUID != "a-b-c-d" {
		t.Fatalf("have %s, want %s", byName.UUID, "a-b-c-d")
	}

	byApplyAt, err := db.BlueprintsByApplyAt(context.Background(), "Enroll")
	if err != nil {
		t.Fatalf("getting blueprint by ApplyAt: %v", err)
	}
	if len(byApplyAt) != 2 {
		t.Fatalf("multiple blueprints not saved correctly")
	}
}

func TestList(t *testing.T) {
	db := setupDB(t)
	bp1 := &blueprint.Blueprint{
		UUID: "a-b-c-d",
		Name: "blueprint-1",
	}
	bp2 := &blueprint.Blueprint{
		UUID: "e-f-g-h",
		Name: "blueprint-2",
	}

	if err := db.Save(bp1); err != nil {
		t.Fatalf("saving blueprint-1 to datastore: %v", err)
	}

	if err := db.Save(bp2); err != nil {
		t.Fatalf("saving blueprint-2 to datastore: %v", err)
	}

	bps, err := db.List()
	if err != nil {
		t.Fatalf("listing blueprints: %v", err)
	}
	if len(bps) != 2 {
		t.Fatalf("expected %d, found %d", 2, len(bps))
	}
}

func TestDelete(t *testing.T) {
	db := setupDB(t)
	bp1 := &blueprint.Blueprint{
		UUID: "a-b-c-d",
		Name: "blueprint",
	}

	if err := db.Save(bp1); err != nil {
		t.Fatalf("saving blueprint to datastore: %v", err)
	}

	if err := db.Delete("blueprint"); err != nil {
		t.Fatalf("deleting blueprint in datastore: %v", err)
	}

	_, err := db.BlueprintByName("blueprint")
	if err == nil {
		t.Fatalf("expected blueprint to be deleted: %v", err)
	}
}

func setupDB(t *testing.T) *DB {
	f, _ := ioutil.TempFile("", "bolt-")
	f.Close()
	os.Remove(f.Name())

	db, err := bolt.Open(f.Name(), 0777, nil)
	if err != nil {
		t.Fatalf("couldn't open bolt, err %v\n", err)
	}
	profileDB, err := profile.NewDB(db)
	if err != nil {
		t.Fatalf("couldn't create profile DB, err %v\n", err)
	}

	blueprintDB, err := NewDB(db, profileDB)
	if err != nil {
		t.Fatalf("couldn't create blueprint DB, err %v\n", err)
	}
	return blueprintDB
}
