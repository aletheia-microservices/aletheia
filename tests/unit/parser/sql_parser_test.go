package parser_test

import (
	"slices"
	"testing"

	appparser "github.com/aletheia-microservices/aletheia/internal/app/parser"
)

func TestParseSQLReadTable(t *testing.T) {
	_, filterFields, tables, ok := appparser.ParseSQLRead("catalogue_db", "SELECT * FROM sock WHERE sock.SockID=?")
	if !ok {
		t.Fatalf("ParseSQLRead failed")
	}
	if !slices.Equal(tables, []string{"sock"}) {
		t.Errorf("tables = %v, want [sock]", tables)
	}
	if len(filterFields) != 1 {
		t.Errorf("filter fields = %v, want 1", filterFields)
	}
}

// the tables of a FROM with JOINs are listed from left to right, so the first one is the main table
// (e.g., the base query of sockshop CatalogueService.List)
func TestParseSQLReadJoinTables(t *testing.T) {
	stmt := `SELECT sock.SockID, sock.Name FROM sock
		JOIN sock_tag allsocktags ON sock.SockID=allsocktags.SockID
		JOIN tag alltags ON allsocktags.TagID=alltags.TagID`

	_, _, tables, ok := appparser.ParseSQLRead("catalogue_db", stmt)
	if !ok {
		t.Fatalf("ParseSQLRead failed")
	}
	if !slices.Equal(tables, []string{"sock", "sock_tag", "tag"}) {
		t.Errorf("tables = %v, want [sock sock_tag tag]", tables)
	}
}
