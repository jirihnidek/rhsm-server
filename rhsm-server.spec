%global goipath         github.com/jirihnidek/rhsm-server
Version:                0.0.1

%gometa -L -f

Name:           rhsm-server
Release:        %autorelease
Summary:        The systemd service providing RHSM Varlink API
License:        Apache-2.0 AND BSD-3-Clause AND GPL-3.0-only AND MIT
URL:            %{gourl}
Source0:        %{gosource}
Source1:        %{archivename}-vendor.tar.bz2
Source2:        go-vendor-tools.toml

BuildRequires:  go-vendor-tools
BuildRequires:  systemd-rpm-macros
BuildRequires:  askalono-cli

%description
The rhsm-server.service provides RHSM Varlink API for client tools.

# --- begin prep ---
%prep
%goprep -A
%setup -q -T -D -a1 %{forgesetupargs}
# ---- end prep ----

%generate_buildrequires
%go_vendor_license_buildrequires -c %{S:2}

# --- begin build ---
%build
export GO_LDFLAGS="-X github.com/jirihnidek/rhsm-server/pkg/version.Version=%{version}"
%gobuild -o %{gobuilddir}/bin/rhsm-server %{goipath}/cmd/rhsm-server
# ---- end build ----

# --- begin install ---
%install
# Licenses
%go_vendor_license_install -c %{S:2}
# Binaries
install -m 0755 -vd                        %{buildroot}%{_libexecdir}/%{name}
install -m 0755 -vp _build/bin/rhsm-server %{buildroot}%{_libexecdir}/%{name}/
# Systemd files
install -m 0755 -vd                     %{buildroot}%{_unitdir}
install -m 0644 -vp data/systemd/rhsm-server.service  %{buildroot}%{_unitdir}/
install -m 0644 -vp data/systemd/rhsm-server.socket   %{buildroot}%{_unitdir}/
install -m 0755 -vd %{buildroot}%{_prefix}/lib/systemd/system-preset/
install -m 0644 -vp data/systemd/presets/50-rhsm-server.preset %{buildroot}%{_prefix}/lib/systemd/system-preset/
# ---- end install ----

# --- begin check ---
%check
%go_vendor_license_check -c %{S:2}
%if 0%{?with_check}
%gocheck2
%endif
# ---- end check ----

# --- begin post ---
%post
%systemd_post rhsm-server.socket
# ---- end post ----

# --- begin preun ---
%preun
%systemd_preun rhsm-server.socket rhsm-server.service
# ---- end preun ----

# --- begin postun ---
%postun
%systemd_postun_with_restart rhsm-server.service
# ---- end postun ----

# --- begin file ---
%files -f %{go_vendor_license_filelist}
# Binaries
%{_libexecdir}/%{name}/rhsm-server
# Systemd
%{_unitdir}/rhsm-server.service
%{_unitdir}/rhsm-server.socket
%{_prefix}/lib/systemd/system-preset/50-rhsm-server.preset
# ---- end files ----

%changelog
%autochangelog
