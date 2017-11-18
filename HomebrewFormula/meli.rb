class Meli < Formula
  desc "Meli is supposed to be a faster alternative to docker-compose.
 Faster in the sense that, Meli will try to pull as many services(docker containers) as it can in parallel."
  homepage "https://github.com/komuW/meli"
  url "https://github.com/komuW/meli/releases/download/v0.1.1.3/meli_0.1.1.3_darwin_amd64.tar.gz"
  version "0.1.1.3"
  sha256 "dfe0f33a3e1f4ed95cc43c7461ef92a68fc7272fda26fa961b73f1d1f4f150ce"

  def install
    bin.install "program"
  end

  test do
    system "#{bin}/program --version"
    ...
  end
end
