package main

import (
	"crypto/md5"
	"fmt"
	"image"
	_ "image/png"
	"net/http"
	"os"
	"testing"

	"github.com/minotar/imgd/storage"
	"github.com/minotar/imgd/storage/memory"
	"github.com/minotar/minecraft"
	"github.com/minotar/minecraft/mockminecraft"
	"github.com/op/go-logging"
	. "github.com/smartystreets/goconvey/convey"
)

const (
	testUser = "clone1018"
)

type SilentWriter struct {
}

func (w SilentWriter) Write(p []byte) (n int, err error) {
	return 0, nil
}

func hashRender(render image.Image) string {
	// md5 hash its pixels
	hasher := md5.New()
	hasher.Write(render.(*image.NRGBA).Pix)
	md5Hash := fmt.Sprintf("%x", hasher.Sum(nil))

	// And return the Hash
	return md5Hash
}

func setupTestCache() {
	if cache == nil {
		cache = make(map[string]storage.Storage)
	}
	cache["cacheUUID"], _ = memory.New(10)
	cache["cacheUserData"], _ = memory.New(10)
	cache["cacheSkin"], _ = memory.New(10)
	cache["cacheSkinTransient"], _ = memory.New(10)
}

func setupTestClient() func() {
	mux := mockminecraft.ReturnMux()
	rt, shutdown := mockminecraft.Setup(mux)

	mcClient = minecraft.NewMinecraft()
	mcClient.Client = &http.Client{Transport: rt}
	mcClient.UsernameAPI.SkinURL = "http://skins.example.net/skins/"
	mcClient.UsernameAPI.CapeURL = "http://skins.example.net/capes/"
	return shutdown
}

func TestMain(m *testing.M) {
	logBackend := logging.NewLogBackend(SilentWriter{}, "", 0)
	stats = MakeStatsCollector()
	setupConfig()
	setupLog(logBackend)
	setupTestCache()
	//setupMcClient()
	shutdown := setupTestClient()
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func TestRenders(t *testing.T) {
	Convey("GetHead should return a valid image", t, func() {
		skin := fetchUsernameSkin(testUser)
		err := skin.GetHead(20)

		So(skin.Processed, ShouldNotBeNil)
		So(err, ShouldBeNil)

		hash := hashRender(skin.Processed)
		So(hash, ShouldEqual, "257972f3af185ce1f92200156d3bfe97")
	})

	Convey("GetHelm should return a valid image", t, func() {
		skin := fetchUsernameSkin(testUser)
		err := skin.GetHelm(20)

		So(skin.Processed, ShouldNotBeNil)
		So(err, ShouldBeNil)

		hash := hashRender(skin.Processed)
		So(hash, ShouldEqual, "257972f3af185ce1f92200156d3bfe97")
	})

	Convey("GetCube should return a valid image", t, func() {
		skin := fetchUsernameSkin(testUser)
		err := skin.GetCube(20)

		So(skin.Processed, ShouldNotBeNil)
		So(err, ShouldBeNil)

		hash := hashRender(skin.Processed)
		So(hash, ShouldEqual, "a253bb68f5ed938eb235ff1e3807940c")
	})

	Convey("GetBust should return a valid image", t, func() {
		skin := fetchUsernameSkin(testUser)
		err := skin.GetBust(20)

		So(skin.Processed, ShouldNotBeNil)
		So(err, ShouldBeNil)

		hash := hashRender(skin.Processed)
		So(hash, ShouldEqual, "1a4ba8567619350883923bde97626ae6")
	})

	Convey("GetBody should return a valid image", t, func() {
		skin := fetchUsernameSkin(testUser)
		err := skin.GetBody(20)

		So(skin.Processed, ShouldNotBeNil)
		So(err, ShouldBeNil)

		hash := hashRender(skin.Processed)
		So(hash, ShouldEqual, "987cc81d031386f429b12da2cd5c1ff2")
	})

	Convey("GetArmorBust should return a valid image", t, func() {
		skin := fetchUsernameSkin(testUser)
		err := skin.GetArmorBust(20)

		So(skin.Processed, ShouldNotBeNil)
		So(err, ShouldBeNil)

		hash := hashRender(skin.Processed)
		So(hash, ShouldEqual, "1a4ba8567619350883923bde97626ae6")
	})

	Convey("GetArmorBody should return a valid image", t, func() {
		skin := fetchUsernameSkin(testUser)
		err := skin.GetArmorBody(20)

		So(skin.Processed, ShouldNotBeNil)
		So(err, ShouldBeNil)

		hash := hashRender(skin.Processed)
		So(hash, ShouldEqual, "987cc81d031386f429b12da2cd5c1ff2")
	})
}

func setupBench() func() {
	logBackend := logging.NewLogBackend(SilentWriter{}, "", 0)
	stats = MakeStatsCollector()
	setupConfig()
	setupLog(logBackend)
	setupTestCache()
	//setupMcClient()
	shutdown := setupTestClient()
	return shutdown
}

func BenchmarkGetHead(b *testing.B) {
	defer setupBench()()
	skin := fetchUsernameSkin(testUser)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		skin.GetHead(20)
	}
}

func BenchmarkGetHelm(b *testing.B) {
	defer setupBench()()
	skin := fetchUsernameSkin(testUser)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		skin.GetHelm(20)
	}
}

func BenchmarkGetCube(b *testing.B) {
	defer setupBench()()
	skin := fetchUsernameSkin(testUser)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		skin.GetCube(20)
	}
}

func BenchmarkGetBust(b *testing.B) {
	defer setupBench()()
	skin := fetchUsernameSkin(testUser)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		skin.GetBust(20)
	}
}

func BenchmarkGetBody(b *testing.B) {
	defer setupBench()()
	skin := fetchUsernameSkin(testUser)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		skin.GetBody(20)
	}
}

func BenchmarkGetArmorBust(b *testing.B) {
	defer setupBench()()
	skin := fetchUsernameSkin(testUser)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		skin.GetArmorBust(20)
	}
}

func BenchmarkGetArmorBody(b *testing.B) {
	defer setupBench()()
	skin := fetchUsernameSkin(testUser)
	b.ResetTimer()

	for n := 0; n < b.N; n++ {
		skin.GetArmorBody(20)
	}
}
