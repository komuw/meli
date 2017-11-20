class Meli < Formula
  desc "Meli is supposed to be a faster alternative to docker-compose.
 Faster in the sense that, Meli will try to pull as many services(docker containers) as it can in parallel."
  homepage "https://github.com/komuW/meli"
  url "https://github.com/komuW/meli/releases/download/v0.1.3/meli_0.1.3_darwin_amd64.tar.gz"
  version "0.1.3"
  sha256 "e6eeb4bcb3852fe90f58b977b31d1ccc5445986644222eac979da3a92fcef333"

  def install
    bin.install "program"
  end

  test do
    system "#{bin}/meli --version"
  end
end
