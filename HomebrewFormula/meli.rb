class Meli < Formula
  desc "Meli is supposed to be a faster alternative to docker-compose.
 Faster in the sense that, Meli will try to pull as many services(docker containers) as it can in parallel."
  homepage "https://github.com/komuW/meli"
  url "https://github.com/komuW/meli/releases/download/v0.1.2.1.1b/meli_0.1.2.1.1b_darwin_amd64.tar.gz"
  version "0.1.2.1.1b"
  sha256 "49a3df1a7d2dcc96f7d4debff32f662089a8764f67f858aad410347ef7fff458"

  def install
    bin.install "program"
  end

  def caveats
    "meli -up"
  end

  test do
    system "#{bin}/program --version"
    ...
  end
end
