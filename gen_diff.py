import re
from string import Template

def get_fields(struct):
	state = re.sub(r"\/\/.*", "", struct)
	fields = re.findall(r"struct\s{(.*)}", state, re.S)[0]
	fields_types = re.findall(r"(\w+)\s+([^\s]+)", fields)
	return fields_types

def get_name(struct):
	return re.findall(r"type\s+(\w+)\s+", struct)[0]

def diff_fields(field, ftype):
	stype = re.findall(r"^(\*?[A-Z]\w+)$", ftype)
	if len(stype) > 0:
		tpl = Template("""
	diff$field := current.$field.Diff(old.$field)
	if len(diff$field) > 0 {
		result["$field"] = diff$field
	}
""")
		return tpl.substitute(field=field)
	if ftype.startswith('map'):
		return diff_map_fields(field, ftype)
	if ftype.startswith('['):
		return diff_arr_fields(field, ftype)
	tpl = Template("""
	if old.$field != current.$field {
		result["$field"] = current.$field
	}
""")
	return tpl.substitute(field=field)

def diff_map_fields(field, ftype):
	subtype = re.findall(r"^(?:map)?\[(.*?)\](.*?\*?(\w+))$", ftype)
	#print(field, ftype, subtype)
	stype = re.findall(r"^(\w+)$", subtype[0][2])
	if subtype[0][1].startswith('map'):
		difftemp = Template("""diff := ${type}MapDiff(val, old$field)
			if len(diff) > 0 {
				result$field[key] = diff
			}""")
	elif subtype[0][1].startswith('['):
		difftemp = Template("""diff := ${type}ArrayDiff(val, old$field)
			if _, ok := diff.(bool); !ok {
				result$field[key] = diff
			}""")
	elif len(re.findall(r"([A-Z]+)",subtype[0][1])) > 0:
		difftemp = Template("""diff := val.Diff(old$field)
			if len(diff) > 0 {
				result$field[key] = diff
			}""")
	else:
		difftemp = Template("""if old$field != val {
				result$field[key] = val
			}""")

	difftpl = difftemp.substitute(type=stype[0], field=field)

	tpl = Template("""
	result$field := map[$keytype]interface{}{}
	for key, val := range current.$field {
		if old$field, ok := old.$field[key]; ok {
			$diff
		} else {
			result$field[key] = val
		}
	}
	for key, _ := range old.$field {
		if _, ok := current.$field[key]; !ok {
			result$field[key] = "_del_"
		}
	}
	if len(result$field) > 0 {
		result["$field"] = result$field
	}
""")
	return tpl.substitute(diff=difftpl, field=field, keytype=subtype[0][0])

def diff_arr_fields(field, ftype):
	subtype = re.findall(r"^(?:map)?\[.*?\](.*?\*?(\w+))$", ftype)
	#print(field, ftype, subtype)
	stype = re.findall(r"^(\w+)$", subtype[0][1])

	if subtype[0][0].startswith('map'):
		difftemp = Template("""diff := ${type}MapDiff(current.$field[index], old.$field[index])
			if len(diff) > 0 {
				result$field[index] = diff""")
	elif subtype[0][0].startswith('['):
		difftemp = Template("""diff := ${type}ArrayDiff(current.$field[index], old.$field[index])
			if _, ok := diff.(bool); !ok {
				result$field[index] = diff""")
	elif len(re.findall(r"([A-Z]+)",subtype[0][0])) > 0:
		difftemp = Template("""diff := current.$field[index].Diff(old.$field[index])
			if len(diff) > 0 {
				result$field[index] = diff""")
	else:
		difftemp = Template("""if old.$field[index] != current.$field[index] {
				result$field[index] = current.$field[index]""")
	
	difftpl = difftemp.substitute(type=stype[0], field=field)

	tpl = Template("""
	result$field := make([]interface{}, len(current.$field))
	delta = 0

	if len(current.$field) == 0 || len(old.$field) == 0 {
		if len(current.$field) == 0 && len(old.$field) != 0 {
			result["$field"] = []interface{}{}
		} else if len(current.$field) != 0 && len(old.$field) == 0 {
			result["$field"] = current.$field
		}
	} else {
		for index, _ := range current.$field {
			if index >= len(old.$field) {
				result$field[index] = current.$field[index]
				delta = delta + 1
				continue
			}
			$diff
				delta = delta + 1
			} else {
				result$field[index] = "__"
			}
		}
		if delta > 0 || len(current.$field) != len(old.$field) {
			result["$field"] = result$field
		}
	}
""")
	return tpl.substitute(diff=difftpl, field=field)

def gen_diff(struct):
	name = get_name(struct)
	result = ""
	#pointer = ""
	rettype = "map[string]"
	delta = ""
	diffs = []
	for field, ftype in get_fields(struct):
		if ftype.startswith('func('):
			return ""
		diffs.append(diff_fields(field, ftype))
		if ftype.startswith('['):
			delta = """delta := 0
"""

	tpl = Template("""func (current $name) Diff(old $name) ${rettype}interface{} {
	result := map[string]interface{}{}
	$delta$body
	${result}return result
}""")

	if name == "State":
		result = """if len(result) == 0 {
		return false
	}
	"""
		#pointer = "*"
		rettype = ""
	return tpl.substitute(name=name, body=''.join(diffs), result=result, rettype=rettype, delta=delta)

def gen_type_array(name):
	pointer = ""
	diff = """if old[index] != current[index] {
				result[index] = current[index]"""
	if len(re.findall(r"([A-Z]+)",name)) > 0:
		pointer = "*"
		diff = """diff := current[index].Diff(old[index])
			if len(diff) > 0 {
				result[index] = diff"""
	
	tpl = Template("""func ${name}ArrayDiff(current, old []$pointer$name) interface{} {
	result := make([]interface{}, len(current))
	delta := 0

	if len(current) == 0 || len(old) == 0 {
		if len(current) == 0 && len(old) != 0 {
			return []interface{}{}
		} else if len(current) != 0 && len(old) == 0 {
			return current
		} else {
			return false
		}
	} else {
		for index, _ := range current {
			if index >= len(old) {
				result[index] = current[index]
				delta = delta + 1
				continue
			}
			$diff
				delta = delta + 1
			} else {
				result[index] = "__"
			}
		}
		if delta == 0 && len(current) == len(old) {
			return false
		}
	}
	return result
}""")
	return tpl.substitute(name=name,pointer=pointer,diff=diff)

#TEST INPUT
#structs = [state, kingupgrade, kingattrs, stonegrade, player, projectile, effect, target, unit]

#cat models/*.go | python gen_diff.py > models/differ.go
import sys
data = ''.join(sys.stdin.readlines())

ignorestructs = ["Ability", "Logic", "AbOptions", "Links", "ProjLink", "Wave"]
structs = []

for struct in re.findall(r"type\s+\w+\s+struct\s+{.*?^}", data, re.S | re.M):
	name = get_name(struct)
	try:
		if name not in ignorestructs:
			structs.append(struct)
		#print(struct)
		#print("\n//" + name)
		#print(gen_diff(struct))
		#print(gen_type_array(name))
	except Exception:
		print('//Error for ' + name)
		continue

#print(structs)

fulldiffer = "package models"
gap = """

"""
pointstructs = ["State"]

for strct in structs:
	fulldiffer =  fulldiffer + gap + gen_diff(strct)

for strct in structs:
	for field, ftype in get_fields(strct):
		stname = re.findall(r"\*([A-Z]\w+)",ftype)
		if len(stname) > 0:
			pointstructs.append(stname[0])
pointstructs = set(pointstructs)

for strct in pointstructs:
	fulldiffer = fulldiffer.replace(strct + ")","*" + strct + ")")

arraydiffs = re.findall(r"\s(\w+)ArrayDiff", fulldiffer)
arraydiffs = set(arraydiffs)

for arrdiff in arraydiffs:
	fulldiffer =  fulldiffer + gap + gen_type_array(arrdiff)

#OUTPUT
print(fulldiffer)
