package service

import "testing"

func TestSetTOMLStringValueReplacesExistingValue(t *testing.T) {
	got := setTOMLStringValue("model = \"gpt-5\"\nopenai_base_url = \"https://api.openai.com/v1\"\n", "openai_base_url", "http://127.0.0.1:43123/v1")
	want := "model = \"gpt-5\"\nopenai_base_url = \"http://127.0.0.1:43123/v1\"\n"
	if got != want {
		t.Fatalf("setTOMLStringValue() = %q, want %q", got, want)
	}
}

func TestSetTOMLStringValueAppendsMissingValue(t *testing.T) {
	got := setTOMLStringValue("model = \"gpt-5\"\n", "openai_base_url", "http://127.0.0.1:43123/v1")
	want := "model = \"gpt-5\"\nopenai_base_url = \"http://127.0.0.1:43123/v1\"\n"
	if got != want {
		t.Fatalf("setTOMLStringValue() = %q, want %q", got, want)
	}
}

func TestSetTOMLStringValueAppendsRootValueBeforeTables(t *testing.T) {
	got := setTOMLStringValue("[projects.\"/repo\"]\ntrust_level = \"trusted\"\n", "model_provider", "gx-openai")
	want := "model_provider = \"gx-openai\"\n\n[projects.\"/repo\"]\ntrust_level = \"trusted\"\n"
	if got != want {
		t.Fatalf("setTOMLStringValue() = %q, want %q", got, want)
	}
}

func TestFindTOMLStringValue(t *testing.T) {
	got := findTOMLStringValue("openai_base_url = \"http://127.0.0.1:43123/v1\"\n", "openai_base_url")
	if got != "http://127.0.0.1:43123/v1" {
		t.Fatalf("findTOMLStringValue() = %q", got)
	}
}

func TestSetTOMLStringValueUpdatesTableValue(t *testing.T) {
	got := setTOMLStringValue("[model_providers.gx-openai]\nname = \"GX\"\nbase_url = \"http://old/v1\"\n", "model_providers.gx-openai.base_url", "http://127.0.0.1:43123/v1")
	want := "[model_providers.gx-openai]\nname = \"GX\"\nbase_url = \"http://127.0.0.1:43123/v1\"\n"
	if got != want {
		t.Fatalf("setTOMLStringValue() = %q, want %q", got, want)
	}
}

func TestRepairCodexConfigTextUsesScopedProvider(t *testing.T) {
	got := repairCodexConfigText("model = \"gpt-5\"\n\n[projects.\"/repo\"]\ntrust_level = \"trusted\"\n")
	want := "model = \"gpt-5\"\nmodel_provider = \"gx-openai\"\n\n[projects.\"/repo\"]\ntrust_level = \"trusted\"\n\n[model_providers.gx-openai]\nname = \"GX OpenAI Proxy\"\nbase_url = \"http://127.0.0.1:43123/v1\"\nwire_api = \"responses\"\nrequires_openai_auth = true\n"
	if got != want {
		t.Fatalf("repairCodexConfigText() = %q, want %q", got, want)
	}
}
