package main

import (
	"regexp"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSpec(t *testing.T) {
	configLoc := "cfg.yml"
	machineExpName := "machine"
	emulatorExpName := "emulator"
	graphicsExpName := "graphics"
	diskExpName := "disk"
	diskCacheExpName := "disk-cache"
	diskIoExpName := "disk-io"

	var content string
	var rules []Expression

	Convey("given config file cfg.yml", t, func() {

		var err error
		Convey(" config file should be loaded ", func() {
			content, err = LoadConfig(configLoc)
			So(err, ShouldBeNil)
			So(content, ShouldNotBeNil)
		})

		Convey("config content should be unmarshaled ", func() {
			rules, err = MarshalConfig(content)
			So(err, ShouldBeNil)
			So(len(rules), ShouldBeGreaterThan, 0)
		})

		Convey("machine expression should be found", func() {
			exp := FindExpressionByName(rules, machineExpName)
			So(exp.Name, ShouldEqual, machineExpName)
		})
		Convey("emulator expression should be found", func() {
			exp := FindExpressionByName(rules, emulatorExpName)
			So(exp.Name, ShouldEqual, emulatorExpName)
		})
		Convey("graphics expression should be found", func() {
			exp := FindExpressionByName(rules, graphicsExpName)
			So(exp.Name, ShouldEqual, graphicsExpName)
		})
		Convey("disk expression should be found", func() {
			exp := FindExpressionByName(rules, diskExpName)
			So(exp.Name, ShouldEqual, diskExpName)
		})
		Convey("disk-cache expression should be found", func() {
			exp := FindExpressionByName(rules, diskCacheExpName)
			So(exp.Name, ShouldEqual, diskCacheExpName)
		})
		Convey("disk-io expression should be found", func() {
			exp := FindExpressionByName(rules, diskIoExpName)
			So(exp.Name, ShouldEqual, diskIoExpName)
		})
	})

	// Only pass t into top-level Convey calls
	Convey("given machine is rhel6.6.0", t, func() {
		machine := `
		<type arch='x86_64' machine='rhel6.6.0'>hvm</type>
		<boot dev='hd'/>
		`
		rule := FindExpressionByName(rules, machineExpName)
		expression := rule.Expression
		replace := rule.Replace

		replacedMachine := `
		<type arch='x86_64' machine='pc'>hvm</type>
		<boot dev='hd'/>
		`

		Convey("machine='rhel[\\S\\s]*?' should mache", func() {
			r, err := regexp.Compile(expression)
			So(err, ShouldBeNil)
			matched := r.MatchString(machine)
			So(matched, ShouldBeTrue)
		})

		Convey("Replace machine='rhel6.6.0' to machine='pc'", func() {
			r := regexp.MustCompile(expression)
			target := r.ReplaceAllString(machine, replace)
			So(target, ShouldEqual, replacedMachine)
		})

	})

	Convey("give emulator location", t, func() {
		emulator := `<devices>
		<emulator>/usr/libexec/qemu-kvm</emulator>
		<disk type='file' device='disk'>
		`

		rule := FindExpressionByName(rules, emulatorExpName)
		expression := rule.Expression
		replace := rule.Replace

		replacedEmulator := `<devices>
		<emulator>/usr/bin/kvm-spice</emulator>
		<disk type='file' device='disk'>
		`

		Convey("Should match <emulator>/usr/libexec/qemu-kvm", func() {
			r, err := regexp.Compile(expression)
			So(err, ShouldBeNil)
			matched := r.MatchString(emulator)
			So(matched, ShouldBeTrue)
		})

		Convey("Should replace to <emulator>/usr/bin/kvm-spice ", func() {
			r := regexp.MustCompile(expression)
			target := r.ReplaceAllString(emulator, replace)
			So(target, ShouldEqual, replacedEmulator)
		})
	})

	Convey("Given graphics ", t, func() {
		graphics := `</channel>
		<graphics type='spice' port='5910' autoport='yes' listen='130.120.2.193'>
		<listen type='address' address='130.120.2.193'/>
	  </graphics>
	  <video>
		`
		rule := FindExpressionByName(rules, graphicsExpName)
		expression := rule.Expression
		replace := rule.Replace

		replacedGraphics := `</channel>
		<graphics type='spice'  autoport='yes' />
	  <video>
		`

		Convey("Should match <graphics[\\S\\s]*?</graphics>", func() {
			r, err := regexp.Compile(expression)
			So(err, ShouldBeNil)
			matched := r.MatchString(graphics)
			So(matched, ShouldBeTrue)
		})

		Convey("Should replace to <graphics type='spice'  autoport='yes' /> ", func() {
			r := regexp.MustCompile(expression)
			target := r.ReplaceAllString(graphics, replace)
			So(target, ShouldEqual, replacedGraphics)
		})
	})

	Convey("Give disk file='/home/VPS/***'", t, func() {
		disks := `<emulator>/usr/libexec/qemu-kvm</emulator>
		<disk type='file' device='disk'>
		  <driver name='qemu' type='raw' cache='writethrough' io='native'/>
		  <source file='/home/VPS/gz006.vda'/>
		  <target dev='vda' bus='virtio'/>
		  <address type='pci' domain='0x0000' bus='0x00' slot='0x06' function='0x0'/>
		</disk>
		<disk type='file' device='disk'>
		  <driver name='qemu' type='raw' cache='writethrough' io='native'/>
		  <source file='/home/VPS/gz006.vdb'/>
		  <target dev='vdb' bus='virtio'/>
		  <address type='pci' domain='0x0000' bus='0x00' slot='0x07' function='0x0'/>
		</disk>
		<controller type='usb' index='0' model='ich9-ehci1'>
		`
		// expression := `<source file='/home/VPS/`
		// replace := `<source file='/tank/kvm-pool/gz-tmp/`
		rule := FindExpressionByName(rules, diskExpName)
		expression := rule.Expression
		replace := rule.Replace

		replacedDisks := `<emulator>/usr/libexec/qemu-kvm</emulator>
		<disk type='file' device='disk'>
		  <driver name='qemu' type='raw' cache='writethrough' io='native'/>
		  <source file='/tank/kvm-pool/gz-tmp/gz006.vda'/>
		  <target dev='vda' bus='virtio'/>
		  <address type='pci' domain='0x0000' bus='0x00' slot='0x06' function='0x0'/>
		</disk>
		<disk type='file' device='disk'>
		  <driver name='qemu' type='raw' cache='writethrough' io='native'/>
		  <source file='/tank/kvm-pool/gz-tmp/gz006.vdb'/>
		  <target dev='vdb' bus='virtio'/>
		  <address type='pci' domain='0x0000' bus='0x00' slot='0x07' function='0x0'/>
		</disk>
		<controller type='usb' index='0' model='ich9-ehci1'>
		`

		Convey("Should match <source file='/home/VPS/", func() {
			r, err := regexp.Compile(expression)
			So(err, ShouldBeNil)
			matched := r.MatchString(disks)
			So(matched, ShouldBeTrue)
		})

		Convey("Should replaced", func() {
			r := regexp.MustCompile(expression)
			target := r.ReplaceAllString(disks, replace)
			So(target, ShouldEqual, replacedDisks)
		})
	})

	Convey("give disk with cache option", t, func() {
		diskCaches := `<emulator>/usr/libexec/qemu-kvm</emulator>
		<disk type='file' device='disk'>
		  <driver name='qemu' type='raw' cache='writethrough' io='native'/>
		  <source file='/home/VPS/gz006.vda'/>
		  <target dev='vda' bus='virtio'/>
		  <address type='pci' domain='0x0000' bus='0x00' slot='0x06' function='0x0'/>
		</disk>
		<disk type='file' device='disk'>
		  <driver name='qemu' type='raw' cache='writethrough' io='native'/>
		  <source file='/home/VPS/gz006.vdb'/>
		  <target dev='vdb' bus='virtio'/>
		  <address type='pci' domain='0x0000' bus='0x00' slot='0x07' function='0x0'/>
		</disk>
		<controller type='usb' index='0' model='ich9-ehci1'>
		`
		rule := FindExpressionByName(rules, diskCacheExpName)
		expression := rule.Expression
		replace := rule.Replace

		replacedDiskCaches := `<emulator>/usr/libexec/qemu-kvm</emulator>
		<disk type='file' device='disk'>
		  <driver name='qemu' type='raw'  io='native'/>
		  <source file='/home/VPS/gz006.vda'/>
		  <target dev='vda' bus='virtio'/>
		  <address type='pci' domain='0x0000' bus='0x00' slot='0x06' function='0x0'/>
		</disk>
		<disk type='file' device='disk'>
		  <driver name='qemu' type='raw'  io='native'/>
		  <source file='/home/VPS/gz006.vdb'/>
		  <target dev='vdb' bus='virtio'/>
		  <address type='pci' domain='0x0000' bus='0x00' slot='0x07' function='0x0'/>
		</disk>
		<controller type='usb' index='0' model='ich9-ehci1'>
		`

		Convey("Should match cache='writethrough'", func() {
			r, err := regexp.Compile(expression)
			So(err, ShouldBeNil)
			matched := r.MatchString(diskCaches)
			So(matched, ShouldBeTrue)
		})

		Convey("Should replaced", func() {
			r := regexp.MustCompile(expression)
			target := r.ReplaceAllString(diskCaches, replace)
			So(target, ShouldEqual, replacedDiskCaches)
		})
	})

	Convey("give disk with io option", t, func() {
		diskIo := `<emulator>/usr/libexec/qemu-kvm</emulator>
		<disk type='file' device='disk'>
		  <driver name='qemu' type='raw'  io='native'/>
		  <source file='/home/VPS/gz006.vda'/>
		  <target dev='vda' bus='virtio'/>
		  <address type='pci' domain='0x0000' bus='0x00' slot='0x06' function='0x0'/>
		</disk>
		<disk type='file' device='disk'>
		  <driver name='qemu' type='raw'  io='native'/>
		  <source file='/home/VPS/gz006.vdb'/>
		  <target dev='vdb' bus='virtio'/>
		  <address type='pci' domain='0x0000' bus='0x00' slot='0x07' function='0x0'/>
		</disk>
		<controller type='usb' index='0' model='ich9-ehci1'>
		`
		rule := FindExpressionByName(rules, diskIoExpName)
		expression := rule.Expression
		replace := rule.Replace

		replacedDiskCaches := `<emulator>/usr/libexec/qemu-kvm</emulator>
		<disk type='file' device='disk'>
		  <driver name='qemu' type='raw'  />
		  <source file='/home/VPS/gz006.vda'/>
		  <target dev='vda' bus='virtio'/>
		  <address type='pci' domain='0x0000' bus='0x00' slot='0x06' function='0x0'/>
		</disk>
		<disk type='file' device='disk'>
		  <driver name='qemu' type='raw'  />
		  <source file='/home/VPS/gz006.vdb'/>
		  <target dev='vdb' bus='virtio'/>
		  <address type='pci' domain='0x0000' bus='0x00' slot='0x07' function='0x0'/>
		</disk>
		<controller type='usb' index='0' model='ich9-ehci1'>
		`

		Convey("Should match io='native'", func() {
			r, err := regexp.Compile(expression)
			So(err, ShouldBeNil)
			matched := r.MatchString(diskIo)
			So(matched, ShouldBeTrue)
		})

		Convey("Should replaced", func() {
			r := regexp.MustCompile(expression)
			target := r.ReplaceAllString(diskIo, replace)
			So(target, ShouldEqual, replacedDiskCaches)
		})
	})

}
