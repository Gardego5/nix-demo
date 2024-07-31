{ pkgs, ... }:
pkgs.dockerTools.pullImage {
  arch = "amd64";
  imageName = "postgres";
  imageDigest =
    "sha256:36ed71227ae36305d26382657c0b96cbaf298427b3f1eaeb10d77a6dea3eec41";
  finalImageName = "postgres";
  finalImageTag = "16.3-alpine3.20";
  sha256 = "9UsGjgW9t1g1e96JMGZDL639AyTng55CouvXZdTdqTI=";
  os = "linux";
}
