class Meli < Formula
  desc "Meli is supposed to be a faster alternative to docker-compose.
 Faster in the sense that, Meli will try to pull as many services(docker containers) as it can in parallel."
  homepage "https://github.com/komuW/meli"
  url "https://github.com/komuW/meli/releases/download/v0.1.7.2/meli_0.1.7.2_darwin_amd64.tar.gz"
  version "0.1.7.2"
  sha256 "a0656db7c600ea6c05447194b0d4b9fc878e0f9cd83adc8f99ec579655f191a1"

  def install
    bin.install "meli"
  end

  test do
    system "#{bin}/meli --version"
  end
end
