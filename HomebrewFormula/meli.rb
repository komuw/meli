class Meli < Formula
  desc "Faster alternative to docker-compose"
  homepage "https://github.com/komuW/meli"
  url "https://github.com/komuW/meli/releases/download/v0.1.8/meli_0.1.8_darwin_amd64.tar.gz"
  version "0.1.8"
  sha256 "dfad89e821509eb12031f5ae0963215a8e4f98f4de25a40ad432965d331a12f8"

  def install
    bin.install "meli"
  end

  test do
    system "#{bin}/meli", "--version"
  end
end
