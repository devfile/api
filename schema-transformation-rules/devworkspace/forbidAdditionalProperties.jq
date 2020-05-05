walk( if type == "object" and .properties.metadata then .properties.metadata.additionalProperties="string" else . end )
walk( if type == "object" and .properties and (.additionalProperties=="string"|not) then .additionalProperties=false else . end )
