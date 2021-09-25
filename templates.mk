MAKEFLAGS += --no-builtin-rules --no-builtin-variables --no-print-directory #--output-sync
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c

# Substitutions
null :=
space := $(null) #
comma := ,
single_quote := '#
single_quote_esc := '"'"'#'
hash := \#
hash_esc := \\\#
dollar_bracket := $${#
dollar_bracket_esc := $$$${#

# Recursive Wildcard Function
define rwildcard
$(foreach d,\
  $(wildcard $(1)*),\
  $(call rwildcard,$(d)/,$(2))\
  $(filter \
    $(subst *,%,$2),\
  $d)\
)
endef

# Find template files recursively and return their target filenames (default pattern: *.*)
# $(call dir_templates,subfolder,*.j2)
define dir_templates
$(patsubst templates/%,%,$(call rwildcard,templates/$(1),$(if $(value 2),$(2),*.*)))
endef

# Render makefile variables by eval-ing the file contents as a variable. Handles lack of newlines,
# escaping double quotes and escaping hashs (#)
# $(eval $(call eval_template,mydir/myfile.txt)) -> mydir/myfile.txt-rendered = <rendered file contents>
define eval_template
$(1)-rendered = $(subst $(hash),$(hash_esc), \
	$(subst $(dollar_bracket),$(dollar_bracket_esc), \
	$(shell awk 1 ORS='~~~' $(1))))
endef

# Render template $(1) to file $(2)
# $(call render_template,template/file1.txt,file1.txt)
define render_template
$(info Creating $(2) from template $(1))
$(eval $(call eval_template,$(1)))
$(shell mkdir -p '$(dir $(2))')
$(shell sed -e $$'s/~~~/\\\n/g' <<< '$(subst $(single_quote),$(single_quote_esc),$($(1)-rendered))' | sed -e '$$d' > $(2)-tmp && mv $(2)-tmp $(2))
endef

# Include a snippet in the template
define include_snippet
$(shell awk 1 ORS='~~~' $(1))
endef
