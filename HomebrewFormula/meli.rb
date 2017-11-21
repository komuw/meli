class Meli < Formula
  desc "Meli is supposed to be a faster alternative to docker-compose.
 Faster in the sense that, Meli will try to pull as many services(docker containers) as it can in parallel."
  homepage "https://github.com/komuW/meli"
  url "https://github.com/komuW/meli/releases/download/v0.1.4/meli_0.1.4_darwin_amd64.tar.gz"
  version "0.1.4"
  sha256 "190deecf8be7bd5101199e173baea1f79a45435287a48b1aa7baf5853d3e4b59"

  def install
    bin.install "meli"
  end

  test do
    system "#{bin}/meli --version"
  end
end
