class Meli < Formula
  desc "Meli is supposed to be a faster alternative to docker-compose.
 Faster in the sense that, Meli will try to pull as many services(docker containers) as it can in parallel."
  homepage "https://github.com/komuW/meli"
  url "https://github.com/komuW/meli/releases/download/v0.1.7.1/meli_0.1.7.1_darwin_amd64.tar.gz"
  version "0.1.7.1"
  sha256 "004ec78867f884476f44caf7a9c794a1c4611cb9adfa2445a77ecc7256c04255"

  def install
    bin.install "meli"
  end

  test do
    system "#{bin}/meli --version"
  end
end
