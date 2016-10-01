package db

import (
	"testing"
)

func TestValidateContainerName(t *testing.T) {
	if validateContainerName("abc") != "abc" {
		t.Error("正しい名前が不要にバリデーションされています")
	}

	if validateContainerName("12345") != "12345" {
		t.Error("正しい名前が不要にバリデーションされています")
	}

	if validateContainerName("ABCDE") != "ABCDE" {
		t.Error("正しい名前が不要にバリデーションされています")
	}

	if validateContainerName("abc123ABC_") != "abc123ABC_" {
		t.Error("正しい名前が不要にバリデーションされています")
	}

	if validateContainerName("*^a%bc#") != "abc" {
		t.Error("バリデーションが正しく行われてません")
	}

	if validateContainerName("a bc d  e") != "abcde" {
		t.Error("バリデーションが正しく行われてません")
	}
}
