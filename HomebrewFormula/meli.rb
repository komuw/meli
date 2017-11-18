class Meli < Formula
  desc "Meli is supposed to be a faster alternative to docker-compose.
 Faster in the sense that, Meli will try to pull as many services(docker containers) as it can in parallel."
  homepage "https://github.com/komuW/meli"
  url "https://github.com/komuW/meli/releases/download/v0.1.1.4/meli_0.1.1.4_darwin_amd64.tar.gz"
  version "0.1.1.4"
  sha256 "8936aa76c2665180140e24ac309cb17caf70c62d211aaa81dc20bf3959e9380b"

  def install
    bin.install "program"
  end

  test do
    system "#{bin}/program --version"
    ...
  end
end
