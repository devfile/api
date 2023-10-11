IDE-targeted variants of the schemas provide the following difference compared to the main schemas:
- They contain additional non-standard `markdownDescription` attributes that are used by IDEs such a VSCode
to provide markdown-rendered documentation hovers.
- They don't contain `default` attributes, since this triggers unwanted addition of defaulted fields during completion in IDEs.