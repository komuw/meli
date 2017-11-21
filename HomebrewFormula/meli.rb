class Meli < Formula
  desc "Meli is supposed to be a faster alternative to docker-compose.
 Faster in the sense that, Meli will try to pull as many services(docker containers) as it can in parallel."
  homepage "https://github.com/komuW/meli"
  url "https://github.com/komuW/meli/releases/download/v0.1.5/meli_0.1.5_darwin_amd64.tar.gz"
  version "0.1.5"
  sha256 "720b00956f5d829285f232ef07ebfa4b4a20aa56163450fefafd96c02e039893"

  def install
    bin.install "meli"
  end

  test do
    system "#{bin}/meli --version"
  end
end
