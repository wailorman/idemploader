#!/usr/bin/env ruby

require "yaml"
require "pp"
require "pathname"
require "json"

DOCS_SRC_PATH = Pathname.new(File.expand_path("../src", __dir__))
DOCS_BUILD_PATH = Pathname.new(File.expand_path("../build", __dir__))

DOCS_BUILD_PATH.mkdir unless DOCS_BUILD_PATH.exist?

docs = YAML.load_file(DOCS_SRC_PATH.join("openapi.yml"))

docs["components"] = {}
Dir[DOCS_SRC_PATH.join("components", "**", "*.yml")].map { |path| Pathname.new(path) }.each do |absolute_path|
  dir, base = absolute_path.relative_path_from(DOCS_SRC_PATH.join("components")).split
  
  cursor = docs["components"]
  dir.to_s.split("/").each do |segment| 
    cursor = cursor[segment] ||= {} 
  end
  
  cursor[base.basename(".yml").to_s] = YAML.load_file(absolute_path)
end

docs["paths"] = {}
Dir[DOCS_SRC_PATH.join("paths", "**", "*.yml")].map { |path| Pathname.new(path) }.each do |absolute_path|
  spec = YAML.load_file(absolute_path)
  next if !spec || spec.empty?

  dir, base = absolute_path.relative_path_from(DOCS_SRC_PATH.join("paths")).split
  dir = "" if dir.to_s == "root"
  
  docs["paths"]["/#{dir.to_s}"] ||= {}
  docs["paths"]["/#{dir.to_s}"][base.basename(".yml").to_s] = spec
end

File.open(DOCS_BUILD_PATH.join("openapi.yml"), "w") { |file| file.write(docs.to_yaml) }
File.open(DOCS_BUILD_PATH.join("openapi.json"), "w") { |file| file.write(JSON.pretty_generate(docs)) }
