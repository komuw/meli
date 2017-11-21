class Meli < Formula
  desc "Meli is supposed to be a faster alternative to docker-compose.
 Faster in the sense that, Meli will try to pull as many services(docker containers) as it can in parallel."
  homepage "https://github.com/komuW/meli"
  url "https://github.com/komuW/meli/releases/download/v0.1.6/meli_0.1.6_darwin_amd64.tar.gz"
  version "0.1.6"
  sha256 "9a50ac92f8f81f2c53de360bf42b26e47ea775e04eec57d93d35a87b445959b2"

  def install
    bin.install "meli"
  end

  test do
    system "#{bin}/meli --version"
  end
end
