source = ["./bin/pillar-box-server"]
bundle_id = "com.0xdeafcafe.pillar-box-server"

apple_id {
  username = "@env:APPLE_ID"
  password = "@env:APPLE_ID_PASSWORD"
  provider = "@env:ASC_PROVIDER"
}

sign {
  application_identity = "Developer ID Application: Alexander Forbes-Reed (6Z49PF6642)"
}

zip {
  output_path = "./bin/PillarBox.zip"
}
